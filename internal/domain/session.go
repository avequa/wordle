// Package domain реализует доменные файлы игры: слово, словарь, сессия, статус, результат, а также покрытие тестами
package domain

import (
	"fiveletters/internal/apperr"
	"fmt"
	"math/rand/v2"
	"slices"
)

// GameSession - хранит состояние одной сессии: загаданное слово, лимит попыток, историю, статус
type GameSession struct {
	answer      string
	dict        []string
	maxAttempts int
	history     []GuessResult
	status      Status
}

// seedToPCG переводит знаковый seed в беззнаковое состояние генератора PCG
func seedToPCG(seed int64) uint64 {
	return uint64(seed) //nolint:gosec
}

// StartGame - проверяет словарь и maxAttempts, выбирает ответ лок.генератором, возвращает новую сессию
func StartGame(dictionary []string, maxAttempts int, seed int64) (GameSession, error) {
	return StartGameWithPool(dictionary, dictionary, maxAttempts, seed)
}

// StartGameWithPool - как StartGame, но ответ выбирается из answerPool (слова нужного уровня и категории),
// а попытки по-прежнему проверяются по всему dictionary
func StartGameWithPool(dictionary, answerPool []string, maxAttempts int, seed int64) (GameSession, error) {
	if maxAttempts <= 0 {
		return GameSession{}, fmt.Errorf(
			"maxAttempts = %d: %w",
			maxAttempts,
			apperr.ErrInvalidAttempts,
		)
	}

	words, err := normalizeDictionary(dictionary)
	if err != nil {
		return GameSession{}, err
	}

	pool, err := normalizeDictionary(answerPool)
	if err != nil {
		return GameSession{}, fmt.Errorf("пул ответов: %w", err)
	}

	rng := rand.New(rand.NewPCG(seedToPCG(seed), 0))

	session := GameSession{
		answer:      pool[rng.IntN(len(pool))],
		dict:        words,
		maxAttempts: maxAttempts,
		history:     make([]GuessResult, 0, maxAttempts),
		status:      StatusInProgress,
	}

	if err := session.checkInvariants(); err != nil {
		return GameSession{}, err
	}

	return session, nil
}

func (s *GameSession) checkInvariants() error {
	if s.status == StatusUnknown || s.answer == "" || s.maxAttempts <= 0 || len(s.dict) == 0 {
		return apperr.ErrSessionNotStart
	}

	if len(s.history) > s.maxAttempts {
		return fmt.Errorf(
			"использовано %d попыток из %d: %w",
			len(s.history),
			s.maxAttempts,
			apperr.ErrSessionInvariant,
		)
	}

	if s.status == StatusInProgress && len(s.history) == s.maxAttempts {
		return fmt.Errorf(
			"игра не завершена при исчерпанных попытках: %w",
			apperr.ErrSessionInvariant,
		)
	}

	if err := validateWord(s.answer); err != nil {
		return fmt.Errorf(
			"%w: ответ сессии: %w",
			apperr.ErrSessionInvariant,
			err,
		)
	}

	if !slices.Contains(s.dict, s.answer) {
		return fmt.Errorf(
			"ответ %q отсутствует в словаре сессии: %w",
			s.answer,
			apperr.ErrSessionInvariant,
		)
	}

	return nil
}

// Answer возвращает загаданное слово
func (s *GameSession) Answer() string { return s.answer }

// Status возвращает текущий статус партии
func (s *GameSession) Status() Status { return s.status }

// MaxAttempts возвращает максимальное число попыток в партии
func (s *GameSession) MaxAttempts() int { return s.maxAttempts }

// UsedAttempts возвращает число израсходованных попыток
func (s *GameSession) UsedAttempts() int { return len(s.history) }

// AttemptsLeft возвращает число оставшихся попыток
func (s *GameSession) AttemptsLeft() int { return s.maxAttempts - len(s.history) }

// History возвращает копию истории проверенных попыток
func (s *GameSession) History() []GuessResult { return slices.Clone(s.history) }
