package std

/**
This file holds the standard "epiphanies" - which are known operations to materialize a larger structure
from its idealized form.  For example, images or archives which need to be materialized into their more
complex structures before others can understand and process their contents.
*/

type epiphanies struct{}

var Epiphanies epiphanies

//func (e epiphanies) JPEG(encoded []byte, model color.Model, decay ...time.Duration) *Epiphany[[]byte, image.Image] {
//	d := atlas.DecayPeriod
//	if len(decay) > 0 {
//		d = decay[0]
//	}
//	epi := NewEpiphany[[]byte, image.Image](func(ideal []byte) (material image.Image, err error) {
//		img, err := jpeg.Decode(bytes.NewReader(ideal))
//		if err != nil {
//			return nil, err
//		}
//		return support.ToColorModel(img), nil
//	}, d)
//	_ = epi.Describe(encoded)
//	return epi
//}
//
//func (e epiphanies) PNG(encoded []byte, model color.Model) (Epiphany[[]byte, image.Image], error) {
//}
