package entities

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"net/url"
	"regexp"
	"strings"
	"url-shortener/internal/domain/exceptions"
)

const (
	shortURLLength = 10
	alphabet       = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
)

var validAliasPattern = regexp.MustCompile(fmt.Sprintf(`^[a-zA-Z0-9_]{%d}$`, shortURLLength))

type URL struct {
	Original string
	Short    string
}

func NewURL(original string) *URL {
	return &URL{
		Original: original,
		Short:    generateShort(),
	}
}

func RestoreURL(original, short string) *URL {
	return &URL{
		Original: original,
		Short:    short,
	}
}

func Validate(alias string) error {
	if !validAliasPattern.MatchString(alias) {
		return exceptions.ErrInvalidChars
	}
	return nil
}

func NormalizeURL(rawURL string) (string, error) {
	parsedURL, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return "", err
	}

	parsedURL.Host = strings.ToLower(parsedURL.Host)
	parsedURL.Scheme = strings.ToLower(parsedURL.Scheme)

	if len(parsedURL.Path) > 1 {
		parsedURL.Path = strings.TrimRight(parsedURL.Path, "/")
	}

	return parsedURL.String(), nil
}

func generateShort() string {
	bytes := make([]byte, shortURLLength)
	for i := range bytes {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(alphabet))))
		bytes[i] = alphabet[num.Int64()]
	}
	return string(bytes)
}
