// Package main implements the Five Letters console game.
package main

import (
	"fiveletters/internal/cli"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	var opts cli.Options

	flag.Usage = func() { _, _ = fmt.Fprint(flag.CommandLine.Output(), usage) }

	flag.Int64Var(&opts.Seed, "seed", time.Now().UnixNano(), "seed для воспроизводимого выбора слова")
	flag.StringVar(&opts.Level, "level", "", "уровень сложности: easy, medium, hard (по умолчанию любой)")
	flag.StringVar(&opts.Category, "category", "", "категория: animals, nature, food, home, city, study, sport")
	flag.BoolVar(&opts.Color, "color", false, "цветной вывод")

	replay := flag.String("replay", "", "повтор партии по --seed и списку попыток через запятую")

	flag.Parse()

	var err error

	if *replay != "" {
		err = cli.Replay(os.Stdout, opts, strings.Split(*replay, ","))
	} else {
		var app *cli.App

		app, err = cli.NewApp(os.Stdin, os.Stdout, opts)
		if err == nil {
			err = app.Run()
		}
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, "ошибка:", err)
		os.Exit(1)
	}
}

// usage - справка по запуску, выводится по -h и при неверном флаге
const usage = `Игра «5 букв»: угадайте слово из пяти русских букв за шесть попыток.
 
Запуск:
  make run ARGS="ФЛАГИ"
  go run ./cmd/fiveletters ФЛАГИ
 
Флаги:
  --seed N               seed для воспроизводимого выбора слова (по умолчанию - текущее время)
  --level УРОВЕНЬ        уровень сложности: easy, medium, hard (по умолчанию - любой)
  --category КАТЕГОРИЯ   категория слов: animals, nature, food, home, city, study, sport
                         (по умолчанию - любая)
  --color                цветной вывод
  --replay СЛОВА         повтор партии: попытки через запятую, вместе с --seed этой партии
  -h, --help             эта справка
 
Примеры:
  make run ARGS="--seed 42"
  make run ARGS="--seed 42 --level easy --category nature --color"
 
Во время партии введите ? вместо слова, чтобы получить подсказку
`
