// Package cli реализует консольный слой игры: чтение ввода, вызов игровой логики и вывод результата
package cli

import (
	"errors"
	"fiveletters/internal/apperr"
	"fiveletters/internal/domain"
	"fmt"
	"io"
	"strings"
)

// Отметки букв в выводе
const (
	markExact   = "✅"
	markPresent = "🟡"
	markAbsent  = "❌"
	markUnknown = "?"
)

// ANSI-коды цветов для цветного вывода
const (
	colorGreen  = "\x1b[32m"
	colorYellow = "\x1b[33m"
	colorGray   = "\x1b[90m"
	colorRed    = "\x1b[31m"
	colorReset  = "\x1b[0m"
)

// Тексты меню
const (
	menuText           = "\nМеню:\n1 - новая игра\n2 - выход\n3 - статистика\n"
	promptMenu         = "Выберите пункт: "
	msgUnknownMenuItem = "Неизвестный пункт меню, введите 1, 2 или 3"
	msgStartGame       = "Игра началась!"
	promptGuess        = "Ваше слово: "
	msgInvalidWord     = "Нужно слово ровно из 5 русских букв"
	msgWordNotInDict   = "Такого слова нет в словаре"
	msgBye             = "До свидания!"

	msgHintHelp = "Введите ? вместо слова, чтобы получить подсказку"
	msgNoStats  = "Статистика: пока не сыграно ни одной партии"
)

// вывод игры
type printer struct {
	wr    io.Writer
	color bool
}

func newPrinter(wr io.Writer, color bool) *printer {
	return &printer{
		wr:    wr,
		color: color,
	}
}

// write - пишет текст в поток вывода
func (p *printer) write(text string) error {
	_, err := io.WriteString(p.wr, text)
	return err
}

// printMenu - печатает меню и приглашение выбрать пункт
func (p *printer) printMenu() error {
	if err := p.write(menuText); err != nil {
		return fmt.Errorf("печать меню: %w", err)
	}

	if err := p.write(promptMenu); err != nil {
		return fmt.Errorf("печать выбора пункта меню: %w", err)
	}

	return nil
}

// printUnknownMenuItem - печатает подсказку при вводе неизвестного пункта меню
func (p *printer) printUnknownMenuItem() error {
	if err := p.write(msgUnknownMenuItem + "\n"); err != nil {
		return fmt.Errorf("печать подсказки меню: %w", err)
	}

	return nil
}

// printGameStart - печатает сообщение о старте игры
func (p *printer) printGameStart(number int, seed int64, opts Options, attempts int) error {
	if err := p.write(msgStartGame); err != nil {
		return fmt.Errorf("старт игры: %w", err)
	}

	text := fmt.Sprintf(
		"\n"+"Сессия №%d (seed %d, уровень: %s, категория: %s). Угадайте слово из %d букв за %d попыток!"+"\n%s\n",
		number,
		seed,
		orAny(opts.Level, "любой"),
		orAny(opts.Category, "любая"),
		domain.WordLength,
		attempts,
		msgHintHelp,
	)

	if err := p.write(text); err != nil {
		return fmt.Errorf("печать сообщения о начале партии: %w", err)
	}

	return nil
}

// printGuessPrompt - печатает приглашение ввести слово
func (p *printer) printWriteGuess() error {
	if err := p.write(promptGuess); err != nil {
		return fmt.Errorf("печать приглашения ввода слова: %w", err)
	}

	return nil
}

// printRejectGuess - печатает подсказку, почему попытка не принята
func (p *printer) printRejectGuess(errFromApplyGuess error) error {
	text := msgInvalidWord
	if errors.Is(errFromApplyGuess, apperr.ErrWordNotInDict) {
		text = msgWordNotInDict
	}

	if err := p.write(text + "\n"); err != nil {
		return fmt.Errorf("печать подсказки к отклоненной попытке: %w", err)
	}

	return nil
}

// printHint - печатает подсказку к загаданному слову
func (p *printer) printHint(hint string) error {
	if err := p.write("Подсказка: " + hint + "\n"); err != nil {
		return fmt.Errorf("печать подсказки: %w", err)
	}

	return nil
}

// printStats - печатает статистику игрока за запуск
func (p *printer) printStats(st stats) error {
	if err := p.write(formatStats(st) + "\n"); err != nil {
		return fmt.Errorf("печать статистики: %w", err)
	}

	return nil
}

// formatStats - статистика игрока за запуск
func formatStats(st stats) string {
	if st.played == 0 {
		return msgNoStats
	}

	return fmt.Sprintf(
		"Статистика: сыграно %d, побед %d (%d%%), поражений %d, текущая серия побед %d, лучшая серия %d",
		st.played,
		st.wins,
		st.wins*100/st.played,
		st.played-st.wins,
		st.streak,
		st.bestStreak,
	)
}

// printBye - печать сообщения-прощания
func (p *printer) printBye() error {
	if err := p.write("\n" + msgBye + "\n"); err != nil {
		return fmt.Errorf("печать прощания: %w", err)
	}

	return nil
}

// formatAttemptsLeft - строка с выводом остатка попыток
func formatAttemptsLeft(left int) string {
	return fmt.Sprintf("Осталось попыток: %d", left)
}

// formatWin - строка победы
func formatWin(used int) string {
	textP := "попытки"

	if used == 1 {
		textP = "попытку"
	}
	if used == 5 || used == 6 {
		textP = "попыток"
	}

	return fmt.Sprintf("Победа! Слово угадано за %d %s", used, textP)
}

// formatLose - строка поражения
func formatLose(answer string) string {
	return "Неудача. Загаданное слово: " + answer
}

// letterMark - эмодзи для отметки одной буквы
func letterMark(letter domain.LetterResult) string {
	switch letter {
	case domain.LetterExact:
		return markExact
	case domain.LetterPresent:
		return markPresent
	case domain.LetterAbsent:
		return markAbsent
	case domain.LetterUnknown:
		return markUnknown
	}

	return markUnknown
}

// printProgress - печатает историю попыток, затем остаток попыток либо итог партии
func (p *printer) printProgress(session *domain.GameSession) error {
	var b strings.Builder

	for _, result := range session.History() {
		b.WriteString(formatGuess(result, p.color) + "\n")
	}

	switch session.Status() {
	case domain.StatusWin:
		b.WriteString(paint(formatWin(session.UsedAttempts()), colorGreen, p.color))

	case domain.StatusLose:
		b.WriteString(paint(formatLose(session.Answer()), colorRed, p.color))

	case domain.StatusInProgress, domain.StatusUnknown:
		b.WriteString(formatAttemptsLeft(session.AttemptsLeft()))
	}

	if err := p.write(b.String() + "\n"); err != nil {
		return fmt.Errorf("печать хода партии: %w", err)
	}

	return nil
}

// formatGuess - строка попытки: «эмодзи + слово из попытки»
func formatGuess(result domain.GuessResult, color bool) string {
	var marks strings.Builder

	for _, letter := range result.Letters {
		marks.WriteString(letterMark(letter))
	}

	marks.WriteString(" ")

	for i, r := range []rune(result.Word) {
		marks.WriteString(paint(string(r), letterColor(result.Letters[i]), color))
	}

	return marks.String()
}

// letterColor - цвет буквы по ее отметке
func letterColor(letter domain.LetterResult) string {
	switch letter {
	case domain.LetterExact:
		return colorGreen
	case domain.LetterPresent:
		return colorYellow
	case domain.LetterAbsent, domain.LetterUnknown:
		return colorGray
	}

	return colorGray
}

// paint - оборачивает текст в ANSI-цвет, если цветной вывод включен
func paint(text, color string, enabled bool) string {
	if !enabled {
		return text
	}

	return color + text + colorReset
}

// orAny - название уровня или категории для вывода; вместо пустого значения - anyText
func orAny(value, anyText string) string {
	if value == "" {
		return anyText
	}

	return value
}
