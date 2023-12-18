package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
)

const EmptyRune = -1

var ErrInvalidString = errors.New("invalid string")

func Unpack(input string) (string, error) {
	if input == "" {
		return input, nil
	}

	var builder strings.Builder
	var prev rune = EmptyRune
	var escaped bool

	for _, ch := range input {
		if !escaped && string(prev) == "\\" {
			if string(ch) == "\\" {
				escaped = true
			}

			prev = ch
			continue
		}

		if count, err := strconv.Atoi(string(ch)); err == nil {
			if prev == EmptyRune {
				return input, ErrInvalidString
			}

			builder.WriteString(strings.Repeat(string(prev), count))
			prev = EmptyRune
		} else {
			if prev != EmptyRune {
				builder.WriteRune(prev)
			}

			prev = ch
		}

		escaped = false
	}

	if prev != EmptyRune {
		builder.WriteRune(prev)
	}

	return builder.String(), nil
}
