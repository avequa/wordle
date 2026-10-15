// Package domain реализует доменные файлы игры: слово, словарь, сессия, статус, результат, а также покрытие тестами
package domain

import (
	"fiveletters/internal/apperr"
	"fmt"
)

// DefaultDictionary - словарь из 85-ти допустимых слов (все слова из Words в том же порядке)
func DefaultDictionary() []string {
	words := make([]string, 0, len(Words()))

	for _, value := range Words() {
		words = append(words, value.Word)
	}

	return words
}

func normalizeDictionary(dict []string) ([]string, error) {
	if len(dict) == 0 {
		return nil, apperr.ErrEmptyDict
	}

	words := make([]string, 0, len(dict))
	seen := make(map[string]struct{}, len(dict))

	for _, value := range dict {
		word, err := checkWord(value)
		if err != nil {
			return nil, fmt.Errorf("ошибка в словаре: %w", err)
		}

		if _, ok := seen[word]; ok {
			return nil, fmt.Errorf(
				"ошибка в словаре: %q: %w",
				word,
				apperr.ErrDuplDict,
			)
		}

		seen[word] = struct{}{}

		words = append(words, word)
	}

	return words, nil
}
