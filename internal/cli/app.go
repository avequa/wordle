// Package cli реализует консольный слой игры: чтение ввода, вызов игровой логики и вывод результата
package cli

import (
	"bufio"
	"errors"
	"fiveletters/internal/apperr"
	"fiveletters/internal/domain"
	"fmt"
	"io"
	"strings"
)

// maxAttempts - макс. число попыток
const maxAttempts = 6

// Пункты главного меню
const (
	menuNewGame = "1"
	menuExit    = "2"
	menuStats   = "3"
)

// hintCommand - ввод вместо слова, чтобы получить подсказку
const hintCommand = "?"

// Options - настройки запуска: seed, уровень сложности, категория слов и цветной вывод
type Options struct {
	Seed     int64
	Level    string // пустая строка - любой уровень
	Category string // пустая строка - любая категория
	Color    bool
}

// App - консольное приложение: читает ввод, вызывает игровую логику, передает результат в printer
type App struct {
	scanner     *bufio.Scanner
	out         *printer
	dict        []string
	gameStarted int

	pool  []string
	opts  Options
	stats stats
	echo  bool // печатать введенное слово (нужно для replay, где ввод не виден)
}

// NewApp - конструктор приложения; возвращает ошибку при неизвестном уровне или категории
func NewApp(
	reader io.Reader,
	out io.Writer,
	opts Options,
) (*App, error) {
	pool, err := domain.AnswerPool(opts.Level, opts.Category)
	if err != nil {
		return nil, fmt.Errorf("настройки игры: %w", err)
	}

	return &App{
		scanner: bufio.NewScanner(reader),
		out:     newPrinter(out, opts.Color),
		dict:    domain.DefaultDictionary(),
		pool:    pool,
		opts:    opts,
	}, nil
}

// Run - главное меню: запускает партии, пока игрок не выберет выход или не закончится ввод
func (a *App) Run() error {
	for {
		if err := a.out.printMenu(); err != nil {
			return fmt.Errorf("меню: %w", err)
		}

		if !a.scanner.Scan() {
			break
		}

		choice := strings.TrimSpace(a.scanner.Text())

		if choice == menuExit {
			break
		}

		if choice == menuStats {
			if err := a.out.printStats(a.stats); err != nil {
				return fmt.Errorf("меню: %w", err)
			}

			continue
		}

		if choice != menuNewGame {
			if err := a.out.printUnknownMenuItem(); err != nil {
				return fmt.Errorf("меню: %w", err)
			}

			continue
		}

		if err := a.playGame(); err != nil {
			return fmt.Errorf("меню: %w", err)
		}
	}

	if err := a.scanner.Err(); err != nil {
		return fmt.Errorf("чтение ввода: %w", err)
	}

	if a.stats.played > 0 {
		if err := a.out.printStats(a.stats); err != nil {
			return fmt.Errorf("выход: %w", err)
		}
	}

	if err := a.out.printBye(); err != nil {
		return fmt.Errorf("выход: %w", err)
	}

	return nil
}

// playGame - одна партия: создает сессию и принимает попытки до победы или поражения
func (a *App) playGame() error {
	seed := a.opts.Seed + int64(a.gameStarted)
	a.gameStarted++

	session, err := domain.StartGameWithPool(a.dict, a.pool, maxAttempts, seed)
	if err != nil {
		return fmt.Errorf("партия №%d: %w", a.gameStarted, err)
	}

	if err := a.out.printGameStart(a.gameStarted, seed, a.opts, session.MaxAttempts()); err != nil {
		return fmt.Errorf("партия №%d: %w", a.gameStarted, err)
	}

	for !session.Status().IsFinished() {
		if err := a.out.printWriteGuess(); err != nil {
			return fmt.Errorf("партия №%d: %w", a.gameStarted, err)
		}

		if !a.scanner.Scan() {
			return nil
		}

		guess := a.scanner.Text()

		if a.echo {
			if err := a.out.write(guess + "\n"); err != nil {
				return fmt.Errorf("партия №%d: %w", a.gameStarted, err)
			}
		}

		if strings.TrimSpace(guess) == hintCommand {
			if err := a.out.printHint(domain.HintFor(session.Answer())); err != nil {
				return fmt.Errorf("партия №%d: %w", a.gameStarted, err)
			}

			continue
		}

		_, guessErr := domain.ApplyGuess(&session, guess)

		var printErr error

		switch {
		case guessErr == nil:
			printErr = a.out.printProgress(&session)
		case errors.Is(guessErr, apperr.ErrInvalidWord), errors.Is(guessErr, apperr.ErrWordNotInDict):
			printErr = a.out.printRejectGuess(guessErr)
		default:
			return fmt.Errorf("партия №%d: %w", a.gameStarted, guessErr)
		}

		if printErr != nil {
			return fmt.Errorf("партия №%d: %w", a.gameStarted, printErr)
		}
	}

	a.stats.add(session.Status() == domain.StatusWin)

	return nil
}

// Replay - повторяет партию по seed и списку попыток: выводит то же, что игрок видел в игре
func Replay(out io.Writer, opts Options, guesses []string) error {
	app, err := NewApp(strings.NewReader(strings.Join(guesses, "\n")+"\n"), out, opts)
	if err != nil {
		return err
	}

	app.echo = true

	return app.playGame()
}
