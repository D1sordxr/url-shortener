package vo

import (
	"github.com/D1sordxr/url-shortener/internal/domain/core/url/errorx"
	"math/rand"
	"time"
)

const (
	upperChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	lowerChars = "abcdefghijklmnopqrstuvwxyz"
	numChars   = "0123456789"
	maxLength  = 128
)

type Alias string

func (a Alias) String() string {
	return string(a)
}

func NewAlias(raw string) (Alias, error) {
	if raw == "" || len(raw) > maxLength {
		return "", errorx.ErrInvalidAliasLength
	}
	return Alias(raw), nil
}

func GenerateAlias() Alias {
	return func(size int) Alias {
		rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

		chars := []rune(upperChars + lowerChars + numChars)

		b := make([]rune, size)
		for i := range b {
			b[i] = chars[rnd.Intn(len(chars))]
		}

		return Alias(b)
	}(rand.Intn(maxLength - 1))
}
