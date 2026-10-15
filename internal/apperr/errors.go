// Package apperr содержит ошибки игры
package apperr

import "errors"

// Ошибки словаря, пользовательского ввода и сессии
var (
	ErrEmptyDict        = errors.New("пустой словарь")
	ErrDuplDict         = errors.New("повторы в словаре")
	ErrInvalidAttempts  = errors.New("некорректное число попыток")
	ErrInvalidWord      = errors.New("некорректное слово")
	ErrSessionNotStart  = errors.New("сессия не начата через StartGame")
	ErrSessionInvariant = errors.New("нарушена целостность текущей сессии")
	ErrWordNotInDict    = errors.New("слова нет в словаре")
	ErrGameFinished     = errors.New("игра завершена")
	ErrUnknownMenuItem  = errors.New("неизвестный пункт меню")

	ErrInputClosed     = errors.New("ввод закончился")
	ErrUnknownLevel    = errors.New("неизвестный уровень сложности")
	ErrUnknownCategory = errors.New("неизвестная категория слов")
)
