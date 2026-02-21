package hw02unpackstring

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	var prevChar rune
	var builder strings.Builder
	var escapedPrevChar bool

	for i, char := range str {
		// Валидация: первый символ не может быть цифрой
		if i == 0 && unicode.IsDigit(char) {
			return "", ErrInvalidString
		}

		// Валидация: две цифры подряд
		if unicode.IsDigit(char) && unicode.IsDigit(prevChar) && !escapedPrevChar {
			return "", ErrInvalidString
		}

		// Валидация: экранируем только цифру или \
		if prevChar == '\\' && !(unicode.IsDigit(char) || char == '\\') {
			return "", ErrInvalidString
		}

		// Игнорируем \ и обрабатываем дальше
		if prevChar == '\\' && !escapedPrevChar {
			prevChar, escapedPrevChar = char, true
			continue
		}

		// Обработка цифры (повторение предыдущего символа)
		if unicode.IsDigit(char) {
			count, err := strconv.Atoi(string(char))
			if err != nil {
				return "", fmt.Errorf("failed to parse count: %w", err)
			}

			// Повторяем предыдущий символ count раз
			builder.WriteString(strings.Repeat(string(prevChar), count))
		} else if prevChar != 0 && !unicode.IsDigit(prevChar) || escapedPrevChar { // Обработка обычного символа
			builder.WriteRune(prevChar)
		}

		prevChar, escapedPrevChar = char, false
	}

	// Добавляем последний символ, если он не цифра
	if prevChar != 0 && !unicode.IsDigit(prevChar) || escapedPrevChar {
		builder.WriteRune(prevChar)
	}

	return builder.String(), nil
}
