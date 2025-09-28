package errorx

import "errors"

var (
	ErrAliasAlreadyExists = errors.New("alias already exists")
)
