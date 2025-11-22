// File: kernels.cl

// Euclidean distance in RGB space
float color_distance(uchar4 c1, uchar4 c2) {
    int dr = (int)c1.x - (int)c2.x;
    int dg = (int)c1.y - (int)c2.y;
    int db = (int)c1.z - (int)c2.z;
    return (float)(dr*dr + dg*dg + db*db);
}

// Find closest palette color for a given pixel
int find_closest_color(
    uchar4 pixel,
    __global const uchar4* palette,
    int palette_size
) {
    float min_distance = FLT_MAX;
    int closest_index = 0;

    for (int i = 0; i < palette_size; i++) {
        float dist = color_distance(pixel, palette[i]);
        if (dist < min_distance) {
            min_distance = dist;
            closest_index = i;
        }
    }

    return closest_index;
}

// Quantize image to palette (no dithering)
__kernel void quantize_to_palette(
    __global const uchar4* input,
    __global uchar4* output,
    __global const uchar4* palette,
    int palette_size,
    int width,
    int height
) {
    int gid = get_global_id(0);

    if (gid >= width * height) return;

    // Get input pixel
    uchar4 pixel = input[gid];

    // Find closest palette color
    int palette_index = find_closest_color(pixel, palette, palette_size);

    // Write quantized pixel
    output[gid] = palette[palette_index];
}

// Floyd-Steinberg dithering kernel
__kernel void floyd_steinberg_dither(
    __global uchar4* image,
    __global const uchar4* palette,
    int palette_size,
    int width,
    int height
) {
    int x = get_global_id(0);
    int y = get_global_id(1);

    if (x >= width || y >= height) return;

    int idx = y * width + x;

    // Get current pixel
    uchar4 old_pixel = image[idx];

    // Find closest palette color
    int palette_idx = find_closest_color(old_pixel, palette, palette_size);
    uchar4 new_pixel = palette[palette_idx];

    // Set quantized pixel
    image[idx] = new_pixel;

    // Calculate quantization error
    int4 error;
    error.x = (int)old_pixel.x - (int)new_pixel.x;
    error.y = (int)old_pixel.y - (int)new_pixel.y;
    error.z = (int)old_pixel.z - (int)new_pixel.z;
    error.w = 0;

    // Distribute error to neighbors (Floyd-Steinberg coefficients)
    // Right pixel (x+1, y): 7/16
    if (x + 1 < width) {
        int idx_r = idx + 1;
        uchar4 pixel_r = image[idx_r];
        pixel_r.x = clamp((int)pixel_r.x + error.x * 7 / 16, 0, 255);
        pixel_r.y = clamp((int)pixel_r.y + error.y * 7 / 16, 0, 255);
        pixel_r.z = clamp((int)pixel_r.z + error.z * 7 / 16, 0, 255);
        image[idx_r] = pixel_r;
    }

    // Bottom-left pixel (x-1, y+1): 3/16
    if (x > 0 && y + 1 < height) {
        int idx_bl = idx + width - 1;
        uchar4 pixel_bl = image[idx_bl];
        pixel_bl.x = clamp((int)pixel_bl.x + error.x * 3 / 16, 0, 255);
        pixel_bl.y = clamp((int)pixel_bl.y + error.y * 3 / 16, 0, 255);
        pixel_bl.z = clamp((int)pixel_bl.z + error.z * 3 / 16, 0, 255);
        image[idx_bl] = pixel_bl;
    }

    // Bottom pixel (x, y+1): 5/16
    if (y + 1 < height) {
        int idx_b = idx + width;
        uchar4 pixel_b = image[idx_b];
        pixel_b.x = clamp((int)pixel_b.x + error.x * 5 / 16, 0, 255);
        pixel_b.y = clamp((int)pixel_b.y + error.y * 5 / 16, 0, 255);
        pixel_b.z = clamp((int)pixel_b.z + error.z * 5 / 16, 0, 255);
        image[idx_b] = pixel_b;
    }

    // Bottom-right pixel (x+1, y+1): 1/16
    if (x + 1 < width && y + 1 < height) {
        int idx_br = idx + width + 1;
        uchar4 pixel_br = image[idx_br];
        pixel_br.x = clamp((int)pixel_br.x + error.x * 1 / 16, 0, 255);
        pixel_br.y = clamp((int)pixel_br.y + error.y * 1 / 16, 0, 255);
        pixel_br.z = clamp((int)pixel_br.z + error.z * 1 / 16, 0, 255);
        image[idx_br] = pixel_br;
    }
}

// Ordered dithering (Bayer matrix)
__kernel void bayer_dither(
    __global const uchar4* input,
    __global uchar4* output,
    __global const uchar4* palette,
    int palette_size,
    int width,
    int height
) {
    // Bayer 8x8 matrix
    __constant int bayer[64] = {
         0, 32,  8, 40,  2, 34, 10, 42,
        48, 16, 56, 24, 50, 18, 58, 26,
        12, 44,  4, 36, 14, 46,  6, 38,
        60, 28, 52, 20, 62, 30, 54, 22,
         3, 35, 11, 43,  1, 33,  9, 41,
        51, 19, 59, 27, 49, 17, 57, 25,
        15, 47,  7, 39, 13, 45,  5, 37,
        63, 31, 55, 23, 61, 29, 53, 21
    };

    int x = get_global_id(0);
    int y = get_global_id(1);

    if (x >= width || y >= height) return;

    int idx = y * width + x;
    uchar4 pixel = input[idx];

    // Get Bayer threshold
    int bayer_x = x % 8;
    int bayer_y = y % 8;
    int threshold = bayer[bayer_y * 8 + bayer_x];

    // Apply threshold to pixel
    int4 adjusted;
    adjusted.x = clamp((int)pixel.x + (threshold - 32), 0, 255);
    adjusted.y = clamp((int)pixel.y + (threshold - 32), 0, 255);
    adjusted.z = clamp((int)pixel.z + (threshold - 32), 0, 255);

    uchar4 adjusted_pixel = (uchar4)(adjusted.x, adjusted.y, adjusted.z, pixel.w);

    // Find closest palette color
    int palette_idx = find_closest_color(adjusted_pixel, palette, palette_size);
    output[idx] = palette[palette_idx];
}