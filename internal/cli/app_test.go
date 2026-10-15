package cli

import (
	"bytes"
	"errors"
	"fiveletters/internal/apperr"
	"fiveletters/internal/domain"
	"io"
	"strings"
	"testing"
)

// testSeed - seed для воспроизводимых сценариев
const testSeed = 42

// runApp - прогоняет приложение на заданном вводе и возвращает весь вывод
func runApp(t *testing.T, input string) string {
	t.Helper()

	return runAppWith(t, Options{Seed: testSeed}, input)
}

// runAppWith - как runApp, но с заданными настройками
func runAppWith(t *testing.T, opts Options, input string) string {
	t.Helper()

	var out bytes.Buffer

	app, err := NewApp(strings.NewReader(input), &out, opts)
	if err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}

	if err := app.Run(); err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}

	return out.String()
}

// answerFor - ответ, который StartGame выберет для seed
func answerFor(t *testing.T, seed int64) string {
	t.Helper()

	session, err := domain.StartGame(domain.DefaultDictionary(), maxAttempts, seed)
	if err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}

	return session.Answer()
}

// wrongGuesses - n слов словаря, не совпадающих с ответом
func wrongGuesses(t *testing.T, answer string, n int) []string {
	t.Helper()

	words := make([]string, 0, n)

	for _, word := range domain.DefaultDictionary() {
		if word != answer && len(words) < n {
			words = append(words, word)
		}
	}

	return words
}

// lines - склеивает строки ввода, завершая каждую переводом строки
func lines(input ...string) string {
	return strings.Join(input, "\n") + "\n"
}

// assertContains - проверяет, что вывод содержит все ожидаемые строки
func assertContains(t *testing.T, out string, want ...string) {
	t.Helper()

	for _, w := range want {
		if !strings.Contains(out, w) {
			t.Errorf("в выводе нет %q:\n%s", w, out)
		}
	}
}

// TestMenuInvalidInput - проверяет, что некорректный пункт меню не роняет программу
func TestMenuInvalidInput(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"текст вместо пункта", lines("abc", "2")},
		{"несуществующий пункт", lines("9", "2")},
		{"пустая строка", lines("", "2")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertContains(t, runApp(t, tt.input), msgUnknownMenuItem, msgBye)
		})
	}
}

// TestPlayGame - проверяет партию: победу, поражение и отклоненные попытки
func TestPlayGame(t *testing.T) {
	answer := answerFor(t, testSeed)
	wrong := wrongGuesses(t, answer, maxAttempts)

	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{
			name:  "победа с первой попытки",
			input: lines("1", answer, "2"),
			want:  []string{"✅✅✅✅✅ " + answer, "Победа! Слово угадано за 1 попытку"},
		},
		{
			name:  "отклоненные попытки не расходуют ход",
			input: lines("1", "небо", "water", "озер0", "абоба", wrong[0], answer, "2"),
			want: []string{
				msgInvalidWord, msgWordNotInDict,
				"Осталось попыток: 5", "Победа! Слово угадано за 2 попытки",
			},
		},
		{
			name:  "шесть ошибок - поражение с показом ответа",
			input: lines(append(append([]string{"1"}, wrong...), "2")...),
			want:  []string{"Осталось попыток: 1", "Неудача. Загаданное слово: " + answer},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertContains(t, runApp(t, tt.input), tt.want...)
		})
	}
}

// TestSeveralGames - проверяет, что без перезапуска можно сыграть несколько партий подряд
func TestSeveralGames(t *testing.T) {
	first := answerFor(t, testSeed)
	second := answerFor(t, testSeed+1)

	out := runApp(t, lines("1", first, "1", second, "2"))

	if got := strings.Count(out, "Победа!"); got != 2 {
		t.Fatalf("побед %d, ожидалось 2:\n%s", got, out)
	}

	assertContains(t, out, "Сессия №1", "Сессия №2", "✅✅✅✅✅ "+second)
}

// TestDeterministicMode - проверяет, что одинаковые seed и ввод дают одинаковый вывод в установленном формате
func TestDeterministicMode(t *testing.T) {
	answer := answerFor(t, testSeed)
	input := lines("1", wrongGuesses(t, answer, 1)[0], answer, "2")

	first := runApp(t, input)

	assertContains(t, first, "Осталось попыток: 5", "✅✅✅✅✅ "+answer, "Победа! Слово угадано за 2 попытки")

	if second := runApp(t, input); second != first {
		t.Fatalf("вывод двух запусков с одним seed различается:\n%s\n---\n%s", first, second)
	}
}

// TestNewAppUnknownOptions - неизвестные уровень или категория дают ожидаемую ошибку
func TestNewAppUnknownOptions(t *testing.T) {
	tests := []struct {
		name    string
		opts    Options
		wantErr error
	}{
		{"неизвестный уровень", Options{Level: "ultra"}, apperr.ErrUnknownLevel},
		{"неизвестная категория", Options{Category: "cars"}, apperr.ErrUnknownCategory},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := NewApp(strings.NewReader(""), io.Discard, tt.opts); !errors.Is(err, tt.wantErr) {
				t.Fatalf("ошибка %v, ожидалась %v", err, tt.wantErr)
			}
		})
	}
}

// TestLevelAndCategory - загаданное слово берется из слов выбранных уровня и категории
func TestLevelAndCategory(t *testing.T) {
	tests := []Options{
		{Seed: testSeed, Level: domain.LevelEasy},
		{Seed: testSeed, Category: domain.CategoryFood},
		{Seed: testSeed, Level: domain.LevelHard, Category: domain.CategorySport},
	}

	for _, opts := range tests {
		t.Run(opts.Level+"/"+opts.Category, func(t *testing.T) {
			pool, err := domain.AnswerPool(opts.Level, opts.Category)
			if err != nil {
				t.Fatalf("неожиданная ошибка: %v", err)
			}

			session, err := domain.StartGameWithPool(domain.DefaultDictionary(), pool, maxAttempts, testSeed)
			if err != nil {
				t.Fatalf("неожиданная ошибка: %v", err)
			}

			out := runAppWith(t, opts, lines("1", session.Answer(), "2"))

			assertContains(t, out, "Победа! Слово угадано за 1 попытку")
		})
	}
}

// TestHint - подсказка печатается и не расходует попытку
func TestHint(t *testing.T) {
	answer := answerFor(t, testSeed)
	out := runApp(t, lines("1", "?", "?", answer, "2"))

	assertContains(t, out, "Подсказка: "+domain.HintFor(answer), "Победа! Слово угадано за 1 попытку")
}

// TestStats - статистика считает только завершенные партии и печатается в меню и при выходе
func TestStats(t *testing.T) {
	first := answerFor(t, testSeed)
	second := answerFor(t, testSeed+1)
	third := answerFor(t, testSeed+2)
	lose := append([]string{"1"}, wrongGuesses(t, second, maxAttempts)...)

	input := lines("3", "1", first, "3")
	input += lines(lose...)
	input += lines("1", third, "1")

	out := runApp(t, input)

	assertContains(t, out,
		msgNoStats,
		"Статистика: сыграно 1, побед 1 (100%), поражений 0, текущая серия побед 1, лучшая серия 1",
		"Статистика: сыграно 3, побед 2 (66%), поражений 1, текущая серия побед 1, лучшая серия 1",
	)
}

// TestColor - при включенном цвете в выводе есть ANSI-коды, без цвета их нет
func TestColor(t *testing.T) {
	answer := answerFor(t, testSeed)
	input := lines("1", answer, "2")

	if out := runAppWith(t, Options{Seed: testSeed, Color: true}, input); !strings.Contains(out, colorGreen) {
		t.Errorf("нет цвета при Color: true:\n%s", out)
	}

	if out := runApp(t, input); strings.Contains(out, "\x1b[") {
		t.Errorf("есть цвет при Color: false:\n%s", out)
	}
}

// TestReplay - повтор по seed и попыткам дает тот же результат, что игра
func TestReplay(t *testing.T) {
	answer := answerFor(t, testSeed)
	wrong := wrongGuesses(t, answer, maxAttempts)

	tests := []struct {
		name    string
		guesses []string
		want    []string
		notWant string
	}{
		{
			name:    "победа, лишние попытки после нее игнорируются",
			guesses: []string{wrong[0], "небо", answer, wrong[1]},
			want: []string{
				"Ваше слово: небо", msgInvalidWord,
				"Осталось попыток: 5", "Победа! Слово угадано за 2 попытки",
			},
			notWant: "Ваше слово: " + wrong[1],
		},
		{
			name:    "поражение",
			guesses: wrong,
			want:    []string{"Неудача. Загаданное слово: " + answer},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var out bytes.Buffer

			if err := Replay(&out, Options{Seed: testSeed}, tt.guesses); err != nil {
				t.Fatalf("неожиданная ошибка: %v", err)
			}

			assertContains(t, out.String(), tt.want...)

			if tt.notWant != "" && strings.Contains(out.String(), tt.notWant) {
				t.Errorf("попытка после конца партии не должна применяться:\n%s", out.String())
			}
		})
	}
}
