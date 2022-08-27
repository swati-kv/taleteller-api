package utils

import gonanoid "github.com/matoous/go-nanoid"

type generator struct {
	alphabet             string
	numbers              string
	alphaNumericCharSize int
	numericCharSize      int
}

func NewGeneratorUtils() IDGeneratorUtils {
	return &generator{
		alphabet:             "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789",
		numbers:              "0123456789",
		alphaNumericCharSize: 10,
		numericCharSize:      4,
	}
}

type IDGeneratorUtils interface {
	GenerateIDWithPrefix(prefix string) (id string, err error)
}

func (g generator) GenerateIDWithPrefix(prefix string) (id string, err error) {
	randomID, err := g.GenerateID()
	if err != nil {
		return
	}

	id = prefix + randomID
	return
}

func (g generator) GenerateID() (finalID string, err error) {
	var alphaNumericChars, numericChars string
	if g.alphaNumericCharSize > 0 {
		alphaNumericChars, err = gonanoid.Generate(g.alphabet, g.alphaNumericCharSize)
		if err != nil {
			return
		}
	}

	if g.numericCharSize > 0 {
		numericChars, err = gonanoid.Generate(g.numbers, g.numericCharSize)
		if err != nil {
			return
		}
	}

	finalID = alphaNumericChars + numericChars
	return
}
