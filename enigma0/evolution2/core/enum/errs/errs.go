package errs

import "errors"

var IdeaCodeRequired = errors.New("a code is required")
var IdeaCodeInvalid = errors.New("the provided code was invalid")
