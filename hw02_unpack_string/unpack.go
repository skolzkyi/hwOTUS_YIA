package hw02unpackstring

import (
	"errors"
	"fmt"
	"strconv"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(input string) (string, error) {
	var res string

	if input == "" {
		return "", nil
	}

	resslice := make([]rune, 0)
	rslice := []rune(input)

	if err := verifRslice(rslice); err != nil {
		return "", err
	}

	for i := 0; i <= len(rslice)-2; i++ {
		if unicode.IsDigit(rslice[i+1]) {
			digit, _ := strconv.Atoi(string(rslice[i+1]))
			for j := 0; j < digit; j++ {
				resslice = append(resslice, rslice[i])
			}
			i++
		} else {
			resslice = append(resslice, rslice[i])
		}
	}

	if !unicode.IsDigit(rslice[len(rslice)-1]) {
		resslice = append(resslice, rslice[len(rslice)-1])
	}

	fmt.Println()
	res = string(resslice)

	return res, nil
}

func verifRslice(input []rune) error {
	for i := 0; i <= len(input)-1; i++ {
		if i == 0 && unicode.IsDigit(input[i]) {
			return ErrInvalidString
		} else if i > 0 && unicode.IsDigit(input[i]) && unicode.IsDigit(input[i-1]) {
			return ErrInvalidString
		}
	}

	return nil
}
