package errorx

import (
	"errors"
)

var (
	ErrURLEmpty             = errors.New("URL cannot be empty")
	ErrURLInvalidFormat     = errors.New("invalid URL format")
	ErrURLMissingScheme     = errors.New("URL must include protocol (http:// or https://)")
	ErrURLMissingHost       = errors.New("URL must include hostname")
	ErrURLUnsupportedScheme = errors.New("only http and https protocols are supported")
	ErrURLParseFailed       = errors.New("failed to parse URL")
)
