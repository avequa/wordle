// Package domain реализует доменные файлы игры: слово, словарь, сессия, статус, результат, а также покрытие тестами
package domain

// LetterResult - проверка (отметка) одной буквы в слове при попытке
type LetterResult int

// Возможные отметки буквы
const (
	LetterUnknown LetterResult = 0
	LetterExact   LetterResult = 1
	LetterPresent LetterResult = 2
	LetterAbsent  LetterResult = 3
)

func (r LetterResult) String() string {
	switch r {
	case LetterExact:
		return "EXACT"
	case LetterPresent:
		return "PRESENT"
	case LetterAbsent:
		return "ABSENT"
	case LetterUnknown:
		return "UNKNOWN"

	default:
		return "UNKNOWN"
	}
}

// GuessResult - результат одной проверенной попытки (само слово, его буквы, статус)
type GuessResult struct {
	Word    string
	Letters [WordLength]LetterResult
	Status  Status
}
