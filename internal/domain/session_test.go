package domain

import (
	"errors"
	"fiveletters/internal/apperr"
	"slices"
	"testing"
)

// newTestSession - создает заведомо корректную сессию для сценариев, где важен не старт игры
func newTestSession(t *testing.T, maxAttempts int, seed int64) GameSession {
	t.Helper()

	session, err := StartGame(DefaultDictionary(), maxAttempts, seed)
	if err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}

	return session
}

func TestStartGameValidation(t *testing.T) {
	tests := []struct {
		name        string
		dict        []string
		maxAttempts int
		wantErr     error
	}{
		{"нулевой словарь", nil, 6, apperr.ErrEmptyDict},
		{"пустой словарь", []string{}, 6, apperr.ErrEmptyDict},
		{"некорректное слово в словаре", []string{"небо"}, 6, apperr.ErrInvalidWord},
		{"дубликат в словаре", []string{"озеро", "ОЗЕРО"}, 6, apperr.ErrDuplDict},
		{"ноль попыток", DefaultDictionary(), 0, apperr.ErrInvalidAttempts},
		{"отрицательное число попыток", DefaultDictionary(), -3, apperr.ErrInvalidAttempts},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			session, err := StartGame(tt.dict, tt.maxAttempts, 42)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("получено %v, ожидалось %v", err, tt.wantErr)
			}

			if session.Status() != StatusUnknown {
				t.Fatalf("при ошибке возвращена начатая сессия: статус %d", session.Status())
			}
		})
	}
}

func TestStartGameInitialState(t *testing.T) {
	dict := DefaultDictionary()

	session, err := StartGame(dict, 6, 42)
	if err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}

	if session.Status() != StatusInProgress {
		t.Errorf("статус %d, ожидался %d", session.Status(), StatusInProgress)
	}

	if session.MaxAttempts() != 6 {
		t.Errorf("максимум попыток %d, ожидалось 6", session.MaxAttempts())
	}

	if session.UsedAttempts() != 0 {
		t.Errorf("использовано попыток %d, ожидалось 0", session.UsedAttempts())
	}

	if session.AttemptsLeft() != session.MaxAttempts() {
		t.Errorf("осталось %d попыток при максимуме %d", session.AttemptsLeft(), session.MaxAttempts())
	}

	if len(session.History()) != 0 {
		t.Errorf("история не пуста: %v", session.History())
	}

	if !slices.Contains(dict, session.Answer()) {
		t.Errorf("ответ %q отсутствует в словаре", session.Answer())
	}
}

func TestStartGameSameSeedSameAnswers(t *testing.T) {
	dict := DefaultDictionary()

	for _, seed := range []int64{0, 42, -7} {
		first := newTestSession(t, 6, seed)

		second, err := StartGame(dict, 6, seed)
		if err != nil {
			t.Fatalf("неожиданная ошибка: %v", err)
		}

		if first.Answer() != second.Answer() {
			t.Fatalf("seed %d: получено %q, ожидалось %q", seed, second.Answer(), first.Answer())
		}
	}
}

func TestCheckInvariants(t *testing.T) {
	t.Run("корректная сессия", func(t *testing.T) {
		session := newTestSession(t, 6, 42)

		if err := session.checkInvariants(); err != nil {
			t.Fatalf("неожиданная ошибка: %v", err)
		}
	})

	t.Run("сессия не создана через StartGame", func(t *testing.T) {
		var session GameSession

		if err := session.checkInvariants(); !errors.Is(err, apperr.ErrSessionNotStart) {
			t.Fatalf("получено %v, ожидалось %v", err, apperr.ErrSessionNotStart)
		}
	})

	t.Run("использовано попыток больше максимума", func(t *testing.T) {
		session := newTestSession(t, 6, 42)
		session.history = make([]GuessResult, session.maxAttempts+1)

		if err := session.checkInvariants(); !errors.Is(err, apperr.ErrSessionInvariant) {
			t.Fatalf("получено %v, ожидалось %v", err, apperr.ErrSessionInvariant)
		}
	})

	t.Run("ответ отсутствует в словаре сессии", func(t *testing.T) {
		session := newTestSession(t, 6, 42)
		session.answer = "мираж"
		session.dict = []string{"озеро", "абоба"}

		if err := session.checkInvariants(); !errors.Is(err, apperr.ErrSessionInvariant) {
			t.Fatalf("получено %v, ожидалось %v", err, apperr.ErrSessionInvariant)
		}
	})
}

// TestStartGameWithPool - ответ берется из пула, а попытки проверяются по всему словарю
func TestStartGameWithPool(t *testing.T) {
	dict := []string{"озеро", "арбуз", "мираж", "школа"}

	tests := []struct {
		name    string
		pool    []string
		wantErr error
	}{
		{"пул из одного слова", []string{"мираж"}, nil},
		{"пул совпадает со словарем", dict, nil},
		{"пустой пул", nil, apperr.ErrEmptyDict},
		{"слова пула нет в словаре", []string{"лимон"}, apperr.ErrSessionInvariant},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			session, err := StartGameWithPool(dict, tt.pool, 6, 42)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("ошибка %v, ожидалась %v", err, tt.wantErr)
			}

			if tt.wantErr == nil && !slices.Contains(tt.pool, session.Answer()) {
				t.Fatalf("ответ %q не из пула %v", session.Answer(), tt.pool)
			}
		})
	}

	t.Run("слово вне пула, но из словаря принимается как попытка", func(t *testing.T) {
		session, err := StartGameWithPool(dict, []string{"мираж"}, 6, 42)
		if err != nil {
			t.Fatalf("неожиданная ошибка: %v", err)
		}

		if _, err := ApplyGuess(&session, "арбуз"); err != nil {
			t.Fatalf("попытка отклонена: %v", err)
		}
	})

	t.Run("StartGame выбирает то же слово, что пул из всего словаря", func(t *testing.T) {
		withPool, _ := StartGameWithPool(dict, dict, 6, 7)
		plain, _ := StartGame(dict, 6, 7)

		if withPool.Answer() != plain.Answer() {
			t.Fatalf("ответы различаются: %q и %q", withPool.Answer(), plain.Answer())
		}
	})
}
