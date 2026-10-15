package domain

import (
	"errors"
	"fiveletters/internal/apperr"
	"testing"
)

func TestNormalizeWord(t *testing.T) {
	tests := []struct {
		name string
		word string
		want string
	}{
		{"верхний регистр", "АБОБА", "абоба"},
		{"смешанный регистр", "АбОбА", "абоба"},
		{"пробелы по краям", "  абоба \t", "абоба"},
		{"буква ё", "абёба", "абеба"},
		{"буква Ё", "АБЁБА", "абеба"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := normalizeWord(tt.word); got != tt.want {
				t.Fatalf("получено %q, ожидалось %q", got, tt.want)
			}
		})
	}
}

func TestValidateWord(t *testing.T) {
	tests := []struct {
		name    string
		word    string
		wantErr error
	}{
		{"корректное слово", "озеро", nil},
		{"четыре буквы", "небо", apperr.ErrInvalidWord},
		{"шесть букв", "корова", apperr.ErrInvalidWord},
		{"пустая строка", "", apperr.ErrInvalidWord},
		{"цифра внутри", "озер0", apperr.ErrInvalidWord},
		{"спецсимвол внутри", "озе-о", apperr.ErrInvalidWord},
		{"english not russian", "water", apperr.ErrInvalidWord},
		{"пробел внутри", "озе о", apperr.ErrInvalidWord},
		{"верхний регистр не проходит без нормализации", "ОЗЕРО", apperr.ErrInvalidWord},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateWord(tt.word); !errors.Is(err, tt.wantErr) {
				t.Fatalf("получено %v, ожидалось %v", err, tt.wantErr)
			}
		})
	}
}

func TestCheckWord(t *testing.T) {
	tests := []struct {
		name    string
		word    string
		want    string
		wantErr bool
	}{
		{"нормализует и принимает", " ОзЁрО ", "озеро", false},
		{"принимает готовое слово", "отбор", "отбор", false},
		{"пробрасывает ошибку валидации", "water", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := checkWord(tt.word)

			if tt.wantErr {
				if !errors.Is(err, apperr.ErrInvalidWord) {
					t.Fatalf("получено %v, ожидалось %v", err, apperr.ErrInvalidWord)
				}

				return
			}

			if err != nil {
				t.Fatalf("неожиданная ошибка: %v", err)
			}

			if got != tt.want {
				t.Fatalf("получено %q, ожидалось %q", got, tt.want)
			}
		})
	}
}
