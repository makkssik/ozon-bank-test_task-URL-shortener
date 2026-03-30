package entities

import "url-shortener/internal/domain/exceptions"

type URL struct {
	Original string
	Short    string
}

func NewURL(original, short string) (*URL, error) {
	if len(short) != 10 {
		return nil, exceptions.ErrInvalidLength
	}
	if !isValidShortChars(short) {
		return nil, exceptions.ErrInvalidChars
	}

	return &URL{Original: original, Short: short}, nil
}

func isValidShortChars(s string) bool {
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' {
			continue
		}
		return false
	}
	return true
}
