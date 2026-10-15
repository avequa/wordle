package domain

import (
	"errors"
	"fiveletters/internal/apperr"
	"slices"
	"testing"
)

// TestWords - проверяет встроенный словарь: слова корректны, у каждого известные уровень и категория и есть подсказка
func TestWords(t *testing.T) {
	for _, info := range Words() {
		t.Run(info.Word, func(t *testing.T) {
			if word, err := checkWord(info.Word); err != nil || word != info.Word {
				t.Errorf("некорректное слово: %v", err)
			}

			if !slices.Contains(Levels(), info.Level) {
				t.Errorf("неизвестный уровень %q", info.Level)
			}

			if !slices.Contains(Categories(), info.Category) {
				t.Errorf("неизвестная категория %q", info.Category)
			}

			if info.Hint == "" {
				t.Error("пустая подсказка")
			}
		})
	}
}

// TestEveryLevelAndCategoryHasWords - в любом сочетании уровня и категории есть что загадать
func TestEveryLevelAndCategoryHasWords(t *testing.T) {
	for _, level := range append(Levels(), "") {
		for _, category := range append(Categories(), "") {
			if _, err := AnswerPool(level, category); err != nil {
				t.Errorf("уровень %q, категория %q: %v", level, category, err)
			}
		}
	}
}

// TestAnswerPool - проверяет выбор слов по уровню и категории
func TestAnswerPool(t *testing.T) {
	tests := []struct {
		name     string
		level    string
		category string
		wantErr  error
	}{
		{"любой уровень и категория", "", "", nil},
		{"только уровень", LevelEasy, "", nil},
		{"только категория", "", CategoryAnimals, nil},
		{"уровень и категория", LevelHard, CategorySport, nil},
		{"неизвестный уровень", "ultra", "", apperr.ErrUnknownLevel},
		{"неизвестная категория", "", "cars", apperr.ErrUnknownCategory},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool, err := AnswerPool(tt.level, tt.category)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("ошибка %v, ожидалась %v", err, tt.wantErr)
			}

			if tt.wantErr != nil {
				return
			}

			for _, word := range pool {
				info := wordInfo(t, word)

				if tt.level != "" && info.Level != tt.level || tt.category != "" && info.Category != tt.category {
					t.Errorf("%q (%s, %s) не подходит под фильтр", word, info.Level, info.Category)
				}
			}
		})
	}

	all, err := AnswerPool("", "")
	if err != nil || !slices.Equal(all, DefaultDictionary()) {
		t.Fatalf("без фильтров должен получиться весь словарь, ошибка: %v", err)
	}
}

// TestHintFor - подсказка есть у слова словаря и пустая у неизвестного слова
func TestHintFor(t *testing.T) {
	if hint := HintFor("озеро"); hint == "" {
		t.Error("нет подсказки к слову «озеро»")
	}

	if hint := HintFor("абоба"); hint != "" {
		t.Errorf("подсказка к неизвестному слову: %q", hint)
	}
}

// wordInfo - находит слово во встроенном словаре
func wordInfo(t *testing.T, word string) WordInfo {
	t.Helper()

	for _, info := range Words() {
		if info.Word == word {
			return info
		}
	}

	t.Fatalf("слова %q нет в словаре", word)

	return WordInfo{}
}
