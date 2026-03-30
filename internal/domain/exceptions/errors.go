package exceptions

import "errors"

var (
	ErrURLNotFound   = errors.New("URL not found")
	ErrAlreadyExists = errors.New("URL already exists")

	ErrInvalidLength = errors.New("short url must be exactly 10 characters long")
	ErrInvalidChars  = errors.New("short url contains invalid characters")
)
