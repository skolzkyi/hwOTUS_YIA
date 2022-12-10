package hw02unpackstring

import (
	"errors"
	"strconv"
	"unicode"
	"unicode/utf8"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(input string) (string, error) {
	var res string

	void, _ := utf8.DecodeRune([]byte(``))

	if input == "" {
		return "", nil
	}

	resslice := make([]rune, 0)
	rslice := []rune(input)

	if err := verifRslice(rslice); err != nil {
		return "", err
	}

	for i := 0; i <= len(rslice)-2; i++ {
		if isDigit(rslice, i+1) {
			digit, _ := strconv.Atoi(string(rslice[i+1]))
			for j := 0; j < digit; j++ {
				resslice = append(resslice, rslice[i])
			}
			i++
		} else if !isSlash(rslice, i) {
			resslice = append(resslice, rslice[i])
			rslice[i] = void
		}
	}

	if !isDigit(rslice, len(rslice)-1) {
		resslice = append(resslice, rslice[len(rslice)-1])
	}

	res = string(resslice)

	return res, nil
}

// не исп. case т.к. условия нельзя свести к одному выражению и по идее
// надо бы возвращать особую ошибку на каждый случай(в тестах можно только одну эту)

//nolint:gocritic
func verifRslice(input []rune) error {
	for i := 0; i <= len(input)-1; i++ {
		if i == 0 && unicode.IsDigit(input[i]) {
			return ErrInvalidString
		} else if i > 0 && isDigit(input, i) && isDigit(input, i-1) {
			return ErrInvalidString
		} else if i > 0 && unicode.IsLetter(input[i]) && isSlash(input, i-1) {
			return ErrInvalidString
		}
	}

	return nil
}

func isDigit(input []rune, index int) bool {
	var res bool
	if unicode.IsDigit(input[index]) && !isSlash(input, index-1) {
		res = true
	}

	return res
}

func isSlash(input []rune, index int) bool {
	var res bool
	slash, _ := utf8.DecodeRune([]byte(`\`))
	if index > 0 {
		if input[index] == slash && input[index-1] != slash {
			res = true
		}
	} else {
		if input[index] == slash {
			res = true
		}
	}

	return res
}
