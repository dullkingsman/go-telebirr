package utils

import (
	"github.com/dullkingsman/go-telebirr/internal/model"
	"math/rand"
	"strings"
	"time"
)

func CreateNonceStr(n int, alphabet ...model.Alphabet) string {
	var _alphabet = model.AlphabetValue.Alphanumerics

	if len(alphabet) > 0 && alphabet[0] != "" {
		_alphabet = alphabet[0]
	}

	var (
		chars = _alphabet.Runes()
		sb    strings.Builder
	)

	for i := 0; i < n; i++ {
		var (
			// since go 1.20 and above have deprecated [rand.Seed]
			seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))
			index      = seededRand.Intn(n) // 0 to n inclusive
		)

		sb.WriteRune(chars[index])
	}

	return sb.String()
}
