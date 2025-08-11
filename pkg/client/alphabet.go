package client

type Alphabet string

var urlFriendlySymbols Alphabet = "-_"
var capitalLetters Alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
var smallLetters Alphabet = "abcdefghijklmnopqrstuvwxyz"
var digits Alphabet = "0123456789"
var letters Alphabet = smallLetters + capitalLetters
var alphanumerics Alphabet = digits + smallLetters + capitalLetters
var all Alphabet = alphanumerics + urlFriendlySymbols

var AlphabetValue = struct {
	UrlFriendlySymbols Alphabet
	CapitalLetters     Alphabet
	SmallLetters       Alphabet
	Digits             Alphabet
	Letters            Alphabet
	Alphanumerics      Alphabet
	All                Alphabet
}{
	UrlFriendlySymbols: urlFriendlySymbols,
	CapitalLetters:     capitalLetters,
	SmallLetters:       smallLetters,
	Digits:             digits,
	Letters:            letters,
	Alphanumerics:      alphanumerics,
	All:                all,
}

var AlphabetValues = []Alphabet{
	urlFriendlySymbols,
	capitalLetters,
	smallLetters,
	digits,
	letters,
	alphanumerics,
	all,
}

func (a Alphabet) String() string {
	return string(a)
}

func (a Alphabet) StringPtr() *string {
	var tmp = a.String()
	return &tmp
}

func (a Alphabet) Runes() []rune {
	return []rune(a.String())
}
