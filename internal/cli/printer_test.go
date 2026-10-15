package cli

import (
	"fiveletters/internal/domain"
	"testing"
)

// TestFormatMessages - проверяет фиксированный формат сообщений из ТЗ
func TestFormatMessages(t *testing.T) {
	arbuz := domain.GuessResult{
		Word: "арбуз",
		Letters: [domain.WordLength]domain.LetterResult{
			domain.LetterAbsent, domain.LetterPresent, domain.LetterAbsent, domain.LetterAbsent, domain.LetterPresent,
		},
	}

	okolo := domain.GuessResult{
		Word: "около",
		Letters: [domain.WordLength]domain.LetterResult{
			domain.LetterExact, domain.LetterAbsent, domain.LetterPresent, domain.LetterAbsent, domain.LetterAbsent,
		},
	}

	tests := []struct {
		name string
		got  string
		want string
	}{
		{"попытка", formatGuess(arbuz, false), "❌🟡❌❌🟡 арбуз"},
		{"повтор букв: отбор и около", formatGuess(okolo, false), "✅❌🟡❌❌ около"},
		{"остаток попыток", formatAttemptsLeft(5), "Осталось попыток: 5"},
		{"победа", formatWin(4), "Победа! Слово угадано за 4 попытки"},
		{"поражение", formatLose("озеро"), "Неудача. Загаданное слово: озеро"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Fatalf("получено %q, ожидалось %q", tt.got, tt.want)
			}
		})
	}
}

// TestFormatGuessColor - буквы слова раскрашиваются по отметкам, эмодзи остаются
func TestFormatGuessColor(t *testing.T) {
	result := domain.GuessResult{
		Word: "около",
		Letters: [domain.WordLength]domain.LetterResult{
			domain.LetterExact, domain.LetterAbsent, domain.LetterPresent, domain.LetterAbsent, domain.LetterAbsent,
		},
	}

	want := "✅❌🟡❌❌ " +
		colorGreen + "о" + colorReset +
		colorGray + "к" + colorReset +
		colorYellow + "о" + colorReset +
		colorGray + "л" + colorReset +
		colorGray + "о" + colorReset

	if got := formatGuess(result, true); got != want {
		t.Fatalf("получено %q, ожидалось %q", got, want)
	}
}

// TestFormatStats - формат статистики за запуск
func TestFormatStats(t *testing.T) {
	tests := []struct {
		name string
		st   stats
		want string
	}{
		{"нет партий", stats{}, msgNoStats},
		{
			"есть партии",
			stats{played: 4, wins: 3, streak: 0, bestStreak: 2},
			"Статистика: сыграно 4, побед 3 (75%), поражений 1, текущая серия побед 0, лучшая серия 2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := formatStats(tt.st); got != tt.want {
				t.Fatalf("получено %q, ожидалось %q", got, tt.want)
			}
		})
	}
}

// TestStatsAdd - серии побед сбрасываются поражением, лучшая серия сохраняется
func TestStatsAdd(t *testing.T) {
	var st stats

	for _, win := range []bool{true, true, false, true} {
		st.add(win)
	}

	want := stats{played: 4, wins: 3, streak: 1, bestStreak: 2}
	if st != want {
		t.Fatalf("получено %+v, ожидалось %+v", st, want)
	}
}
