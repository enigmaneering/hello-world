package errs

import "errors"

var InvalidCode = errors.New("the provided code was invalid")
var InvalidPath = errors.New("the provided path was invalid")
