// Package domain реализует доменные файлы игры: слово, словарь, сессия, статус, результат, а также покрытие тестами
package domain

import (
	"fiveletters/internal/apperr"
	"fmt"
	"strings"
)

// WordLength - фиксированная длина каждого слова в игре
const WordLength = 5

func normalizeWord(word string) string {
	word = strings.ToLower(strings.TrimSpace(word))
	word = strings.ReplaceAll(word, "ё", "е")

	return word
}

func validateWord(word string) error {
	runeWord := []rune(word)

	if len(runeWord) != WordLength {
		return fmt.Errorf(
			"%q: %w",
			word,
			apperr.ErrInvalidWord,
		)
	}

	for _, r := range runeWord {
		if r < 'а' || r > 'я' {
			return fmt.Errorf(
				"%q: %w",
				word,
				apperr.ErrInvalidWord,
			)
		}
	}

	return nil
}

func checkWord(word string) (string, error) {
	normWord := normalizeWord(word)

	if err := validateWord(normWord); err != nil {
		return "", err
	}

	return normWord, nil
}
