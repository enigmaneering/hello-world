package errs

import "errors"

var CodeInvalid = errors.New("the provided code was invalid")
var InvalidPath = errors.New("the provided path was invalid")
