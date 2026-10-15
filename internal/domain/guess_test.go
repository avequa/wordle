package domain

import (
	"errors"
	"fiveletters/internal/apperr"
	"testing"
)

// newSessionWithAnswer - создает сессию с нужным ответом для сценариев, где важен не выбор слова
func newSessionWithAnswer(t *testing.T, answer string, maxAttempts int) GameSession {
	t.Helper()

	session := newTestSession(t, maxAttempts, 67)
	session.answer = answer

	if err := session.checkInvariants(); err != nil {
		t.Fatalf("некорректная тестовая сессия: %v", err)
	}

	return session
}

// TestCompareWords - проверяет расстановку отметок, в том числе при повторах букв
func TestCompareWords(t *testing.T) {
	tests := []struct {
		name   string
		answer string
		guess  string
		want   [WordLength]LetterResult
	}{
		{
			"буква на месте, буква не на месте и лишний повтор",
			"отбор", "около",
			[WordLength]LetterResult{LetterExact, LetterAbsent, LetterPresent, LetterAbsent, LetterAbsent},
		},
		{
			"буквы есть в ответе, но не на своих местах",
			"озеро", "арбуз",
			[WordLength]LetterResult{LetterAbsent, LetterPresent, LetterAbsent, LetterAbsent, LetterPresent},
		},
		{
			"полное совпадение",
			"озеро", "озеро",
			[WordLength]LetterResult{LetterExact, LetterExact, LetterExact, LetterExact, LetterExact},
		},
		{
			"ни одной общей буквы",
			"отбор", "магия",
			[WordLength]LetterResult{LetterAbsent, LetterAbsent, LetterAbsent, LetterAbsent, LetterAbsent},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := compareWords(tt.answer, tt.guess); got != tt.want {
				t.Fatalf("получено %v, ожидалось %v", got, tt.want)
			}
		})
	}
}

// TestApplyGuessRejectsGuess - проверяет отклонение попытки: сессия и число ходов не меняются
func TestApplyGuessRejectsGuess(t *testing.T) {
	tests := []struct {
		name    string
		guess   string
		wantErr error
	}{
		{"четыре буквы", "небо", apperr.ErrInvalidWord},
		{"шесть букв", "корова", apperr.ErrInvalidWord},
		{"цифра внутри", "озер0", apperr.ErrInvalidWord},
		{"спецсимвол внутри", "озе-о", apperr.ErrInvalidWord},
		{"латиница", "water", apperr.ErrInvalidWord},
		{"слова нет в словаре", "абоба", apperr.ErrWordNotInDict},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			session := newSessionWithAnswer(t, "отбор", 6)

			result, err := ApplyGuess(&session, tt.guess)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("получено %v, ожидалось %v", err, tt.wantErr)
			}

			if result != (GuessResult{}) {
				t.Errorf("при ошибке возвращен результат %v", result)
			}

			if session.UsedAttempts() != 0 {
				t.Errorf("использовано попыток %d, ожидалось 0", session.UsedAttempts())
			}

			if len(session.History()) != 0 {
				t.Errorf("история не пуста: %v", session.History())
			}

			if session.Status() != StatusInProgress {
				t.Errorf("статус %d, ожидался %d", session.Status(), StatusInProgress)
			}
		})
	}

	t.Run("сессия не начата", func(t *testing.T) {
		var empty GameSession

		for _, session := range []*GameSession{nil, &empty} {
			result, err := ApplyGuess(session, "около")
			if !errors.Is(err, apperr.ErrSessionNotStart) {
				t.Errorf("получено %v, ожидалось %v", err, apperr.ErrSessionNotStart)
			}

			if result != (GuessResult{}) {
				t.Errorf("при ошибке возвращен результат %v", result)
			}
		}
	})

	t.Run("нарушены инварианты сессии", func(t *testing.T) {
		session := newSessionWithAnswer(t, "отбор", 6)
		session.history = make([]GuessResult, session.maxAttempts+1)

		result, err := ApplyGuess(&session, "около")
		if !errors.Is(err, apperr.ErrSessionInvariant) {
			t.Fatalf("получено %v, ожидалось %v", err, apperr.ErrSessionInvariant)
		}

		if result != (GuessResult{}) {
			t.Errorf("при ошибке возвращен результат %v", result)
		}
	})
}

// TestApplyGuessIgnoresCase - проверяет, что регистр введенного слова не влияет на результат
func TestApplyGuessIgnoresCase(t *testing.T) {
	session := newSessionWithAnswer(t, "отбор", 6)
	want := [WordLength]LetterResult{LetterExact, LetterAbsent, LetterPresent, LetterAbsent, LetterAbsent}

	result, err := ApplyGuess(&session, "ОкОлО")
	if err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}

	if result.Word != "около" {
		t.Errorf("получено %q, ожидалось %q", result.Word, "около")
	}

	if result.Letters != want {
		t.Errorf("получено %v, ожидалось %v", result.Letters, want)
	}

	if result.Status != StatusInProgress {
		t.Errorf("статус результата %d, ожидался %d", result.Status, StatusInProgress)
	}
}

// TestApplyGuessWin - проверяет, что правильный ответ завершает игру победой
func TestApplyGuessWin(t *testing.T) {
	session := newSessionWithAnswer(t, "отбор", 6)
	want := [WordLength]LetterResult{LetterExact, LetterExact, LetterExact, LetterExact, LetterExact}

	result, err := ApplyGuess(&session, "отбор")
	if err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}

	if result.Letters != want {
		t.Errorf("получено %v, ожидалось %v", result.Letters, want)
	}

	if result.Status != StatusWin {
		t.Errorf("статус результата %d, ожидался %d", result.Status, StatusWin)
	}

	if session.Status() != StatusWin {
		t.Errorf("статус сессии %d, ожидался %d", session.Status(), StatusWin)
	}

	if session.UsedAttempts() != 1 {
		t.Errorf("использовано попыток %d, ожидалось 1", session.UsedAttempts())
	}
}

// TestApplyGuessLose - проверяет, что последняя неверная попытка завершает игру поражением
func TestApplyGuessLose(t *testing.T) {
	wrong := []string{"арбуз", "баран", "белка", "берег", "билет", "буква"}
	session := newSessionWithAnswer(t, "отбор", len(wrong))

	var last GuessResult

	for _, guess := range wrong {
		result, err := ApplyGuess(&session, guess)
		if err != nil {
			t.Fatalf("попытка %q: неожиданная ошибка: %v", guess, err)
		}

		last = result
	}

	if last.Status != StatusLose {
		t.Errorf("статус результата %d, ожидался %d", last.Status, StatusLose)
	}

	if session.Status() != StatusLose {
		t.Errorf("статус сессии %d, ожидался %d", session.Status(), StatusLose)
	}

	if session.AttemptsLeft() != 0 {
		t.Errorf("осталось попыток %d, ожидалось 0", session.AttemptsLeft())
	}

	if session.Answer() != "отбор" {
		t.Errorf("ответ %q, ожидался %q", session.Answer(), "отбор")
	}
}

// TestApplyGuessFinishedGame - проверяет, что завершенная игра не принимает новые попытки
func TestApplyGuessFinishedGame(t *testing.T) {
	tests := []struct {
		name  string
		first string
	}{
		{"после победы", "отбор"},
		{"после поражения", "арбуз"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			session := newSessionWithAnswer(t, "отбор", 1)

			first, err := ApplyGuess(&session, tt.first)
			if err != nil {
				t.Fatalf("неожиданная ошибка: %v", err)
			}

			if !first.Status.IsFinished() {
				t.Fatalf("статус результата %d, ожидалась завершенная игра", first.Status)
			}

			result, err := ApplyGuess(&session, "около")
			if !errors.Is(err, apperr.ErrGameFinished) {
				t.Fatalf("получено %v, ожидалось %v", err, apperr.ErrGameFinished)
			}

			if result != (GuessResult{}) {
				t.Errorf("при ошибке возвращен результат %v", result)
			}

			if session.UsedAttempts() != 1 {
				t.Errorf("использовано попыток %d, ожидалось 1", session.UsedAttempts())
			}
		})
	}
}
