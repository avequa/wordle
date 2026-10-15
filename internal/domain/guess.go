// Package domain реализует доменные файлы игры: слово, словарь, сессия, статус, результат, а также покрытие тестами
package domain

import (
	"fiveletters/internal/apperr"
	"fmt"
	"slices"
)

// ApplyGuess - проверяет правила попытки (func checkWord), дописывает ее в историю сессии, обновляет статус:
// победа при полном совпадении, поражение после последней допустимой попытки;
// отклоненная попытка не тратит попытку.
// Игровая модель не зависит от консольного ввода и вывода.
func ApplyGuess(session *GameSession, guess string) (GuessResult, error) {
	if session == nil {
		return GuessResult{}, apperr.ErrSessionNotStart
	}

	if err := session.checkInvariants(); err != nil {
		return GuessResult{}, err
	}

	if session.status.IsFinished() {
		return GuessResult{}, fmt.Errorf(
			"попытка: %q: %w",
			guess,
			apperr.ErrGameFinished,
		)
	}

	word, err := checkWord(guess)
	if err != nil {
		return GuessResult{}, fmt.Errorf("попытка: %w", err)
	}

	if !slices.Contains(session.dict, word) {
		return GuessResult{}, fmt.Errorf(
			"попытка: %q: %w",
			word,
			apperr.ErrWordNotInDict,
		)
	}

	letters := compareWords(session.answer, word)

	result := GuessResult{
		Word:    word,
		Letters: letters,
		Status:  setStatus(letters, len(session.history)+1, session.maxAttempts),
	}

	session.history = append(slices.Clip(session.history), result)
	session.status = result.Status

	return result, nil
}

// compareWords - размечает буквы попытки относительно ответа,
// сначала находятся точные совпадения, затем оставшимися совпадения не на своем месте
func compareWords(answer, guess string) [WordLength]LetterResult {
	var letters [WordLength]LetterResult

	answerRunes := []rune(answer)
	guessRunes := []rune(guess)

	unusedLetters := make(map[rune]int, WordLength)

	// первый проход
	for i, answerRune := range answerRunes {
		if guessRunes[i] == answerRune {
			letters[i] = LetterExact

			continue
		}
		unusedLetters[answerRune]++
	}

	// второй проход
	for i, guessRune := range guessRunes {
		if letters[i] == LetterExact {
			continue
		}

		// количество оставшихся букв, равных текущей букве попытки, больше нуля?
		if unusedLetters[guessRune] > 0 {
			letters[i] = LetterPresent
			unusedLetters[guessRune]--

			continue
		}

		// нет такой буквы в оставшейся части слова
		letters[i] = LetterAbsent
	}

	return letters
}

func isWin(letters [WordLength]LetterResult) bool {
	for _, letter := range letters {
		if letter != LetterExact {
			return false
		}
	}

	return true
}

func setStatus(letters [WordLength]LetterResult, usedAttempts, maxAttempts int) Status {
	switch {
	case isWin(letters):
		return StatusWin

	case usedAttempts >= maxAttempts:
		return StatusLose

	default:
		return StatusInProgress
	}
}
