package exceptions

import "errors"

var (
	ErrAliasCollision = errors.New("alias collision")
	ErrInvalidChars   = errors.New("short url contains invalid characters")
)
