package client

import (
	"crypto"
	cryptoRand "crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"math/rand"
	"reflect"
	"slices"
	"strings"
	"time"
)

type SignatureData []SignatureDataUnit

func NewSignatureData() *SignatureData {
	var tmp = make(SignatureData, 0)
	return &tmp
}

func (s *SignatureData) Add(data interface{}, exclusions SignatureDataExclusions) *SignatureData {
	if s == nil {
		return nil
	}

	*s = append(*s, NewSignDataUnit(data, exclusions))
	return s
}

func (s *SignatureData) AddAll(data []interface{}, exclusions ...SignatureDataExclusions) *SignatureData {
	if s == nil {
		return nil
	}

	for i, d := range data {
		var exclusion SignatureDataExclusions = nil

		if len(exclusions) > i {
			exclusion = exclusions[i]
		}

		*s = append(*s, NewSignDataUnit(d, exclusion))
	}

	return s
}

func (s *SignatureData) Construct(separator ...string) SignatureString {
	if s == nil {
		return ""
	}

	var (
		fieldMap = make(map[string]string)
		fields   []string
	)

	for _, unit := range *s {
		val := reflect.ValueOf(unit.Data)
		typ := reflect.TypeOf(unit.Data)

		for i := 0; i < val.NumField(); i++ {
			field := typ.Field(i)
			jsonKey := field.Tag.Get("json")

			if unit.Exclusions != nil && unit.Exclusions[jsonKey] {
				continue
			}

			if val.Field(i).IsZero() {
				continue
			}

			var value = val.Field(i).String()

			fieldMap[jsonKey] = value
			fields = append(fields, jsonKey)
		}
	}

	return NewSignString(fieldMap, fields, separator...)
}

type SignatureDataExclusions map[string]bool

type SignatureDataUnit struct {
	Data       interface{}
	Exclusions SignatureDataExclusions
}

func NewSignDataUnit(data interface{}, exclusions SignatureDataExclusions) SignatureDataUnit {
	var _exclusions = make(SignatureDataExclusions)

	if exclusions != nil {
		_exclusions = exclusions
	}

	return SignatureDataUnit{
		Data:       data,
		Exclusions: _exclusions,
	}
}

type SignatureString string

func NewSignString(data map[string]string, fields []string, separator ...string) SignatureString {
	var (
		signStrList []string
		seen        = make(map[string]bool)
		_separator  = "&"
	)

	if len(separator) > 0 && separator[0] != "" {
		_separator = separator[0]
	}

	slices.Sort(fields) // just to make sure

	for _, key := range fields {
		if seen[key] { // this is to avoid duplicate keys - fucking headache
			continue
		}

		seen[key] = true
		signStrList = append(signStrList, fmt.Sprintf("%s=%s", key, data[key]))
	}

	return SignatureString(strings.Join(signStrList, _separator))
}

func (s SignatureString) String() string {
	return string(s)
}

func (s SignatureString) Sign(privateKey *rsa.PrivateKey) (string, error) {
	var (
		hashed         = sha256.Sum256([]byte(s.String()))
		signature, err = rsa.SignPSS(cryptoRand.Reader, privateKey, crypto.SHA256, hashed[:], &rsa.PSSOptions{SaltLength: rsa.PSSSaltLengthEqualsHash})
	)

	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(signature), nil
}

func (s SignatureString) VerifySignature(sigBytes string, publicKey *rsa.PublicKey) error {
	var (
		hashed = sha256.Sum256([]byte(s.String()))
		err    = rsa.VerifyPSS(publicKey, crypto.SHA256, hashed[:], []byte(sigBytes), &rsa.PSSOptions{
			SaltLength: rsa.PSSSaltLengthEqualsHash,
			Hash:       crypto.SHA256,
		})
	)

	if err != nil {
		return err
	}

	return nil
}

func CreateNonceStr(n int, alphabet ...Alphabet) string {
	var _alphabet = AlphabetValue.Alphanumerics

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
