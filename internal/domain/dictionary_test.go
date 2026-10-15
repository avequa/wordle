package domain

import (
	"errors"
	"fiveletters/internal/apperr"
	"slices"
	"testing"
)

// minDictSize - минимальный размер словаря
const minDictSize = 50

func TestDefaultDictionary(t *testing.T) {
	t.Run("не менее 50 уникальных слов", func(t *testing.T) {
		dict := DefaultDictionary()
		seen := make(map[string]struct{}, len(dict))

		for _, word := range dict {
			if _, ok := seen[word]; ok {
				t.Fatalf("слово %q встречается дважды", word)
			}

			seen[word] = struct{}{}
		}

		if len(seen) < minDictSize {
			t.Fatalf("уникальных слов %d, ожидалось не менее %d", len(seen), minDictSize)
		}
	})

	t.Run("каждое слово из пяти русских букв", func(t *testing.T) {
		for _, word := range DefaultDictionary() {
			if lenWord := len([]rune(word)); lenWord != WordLength {
				t.Errorf("слово %q: букв %d, ожидалось %d", word, lenWord, WordLength)
			}

			if err := validateWord(word); err != nil {
				t.Errorf("слово %q: %v", word, err)
			}
		}
	})

	t.Run("возвращает независимую копию словаря", func(t *testing.T) {
		dict := DefaultDictionary()
		dict[0] = "абоба"

		if slices.Contains(DefaultDictionary(), "абоба") {
			t.Fatal("изменение результата повлияло на дефолтный (встроенный) словарь")
		}
	})
}

func TestNormalizeDictionary(t *testing.T) {
	tests := []struct {
		name    string
		dict    []string
		want    []string
		wantErr error
	}{
		{"нулевой словарь", nil, nil, apperr.ErrEmptyDict},
		{"пустой словарь", []string{}, nil, apperr.ErrEmptyDict},
		{"короткое слово", []string{"озеро", "небо"}, nil, apperr.ErrInvalidWord},
		{"длинное слово", []string{"корабль"}, nil, apperr.ErrInvalidWord},
		{"цифра в слове", []string{"озер0"}, nil, apperr.ErrInvalidWord},
		{"спецсимвол в слове", []string{"озе-о"}, nil, apperr.ErrInvalidWord},
		{"латиница", []string{"water"}, nil, apperr.ErrInvalidWord},
		{"явный дубликат", []string{"озеро", "озеро"}, nil, apperr.ErrDuplDict},
		{"дубликат после нормализации", []string{"озеро", " ОзЁрО "}, nil, apperr.ErrDuplDict},
		{"нормализует слова и хранит порядок", []string{" ОтБоР ", "ЛАРЁК", "арбуз"}, []string{"отбор", "ларек", "арбуз"}, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := normalizeDictionary(tt.dict)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("получено %v, ожидалось %v", err, tt.wantErr)
				}

				if got != nil {
					t.Fatalf("при ошибке возвращен словарь %v", got)
				}

				return
			}

			if err != nil {
				t.Fatalf("неожиданная ошибка: %v", err)
			}

			if !slices.Equal(got, tt.want) {
				t.Fatalf("получено %v, ожидалось %v", got, tt.want)
			}
		})
	}
}
