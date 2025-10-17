package std

// Pathable represents any type which can return a Path.
type Pathable interface {
	Path() Path
}
