package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/sunr3d/quorum-grep/internal/client"
	"github.com/sunr3d/quorum-grep/internal/config"
	"github.com/sunr3d/quorum-grep/models"
)

func main() {
	flags, err := parseFlags()
	if err != nil {
		fmt.Fprintf(os.Stderr, "parseFlags: %v\n", err)
		fmt.Fprintln(os.Stderr, "Использование утилиты: grep [флаги] шаблон [файлы...]")
		os.Exit(1)
	}

	cfg, err := config.GetConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "config.GetConfig: %v\n", err)
		os.Exit(1)
	}

	client := client.New(cfg)

	for _, file := range flags.Files {
		if err := client.ProcessFile(file, flags.Options); err != nil {
			fmt.Fprintf(os.Stderr, "client.ProcessFile: %v\n", err)
		}
	}
}

// parseFlags - парсит флаги командной строки.
func parseFlags() (*models.GrepConfig, error) {
	opts := models.GrepOptions{}

	flag.IntVar(&opts.After, "A", 0, "напечатать +N строк после найденной строки")
	flag.IntVar(&opts.Before, "B", 0, "напечатать +N строк перед найденной строкой")
	flag.IntVar(&opts.Around, "C", 0, "напечатать +N строк вокруг найденной строки")
	flag.BoolVar(&opts.Count, "c", false, "напечатать только количество найденных строк")
	flag.BoolVar(&opts.IgnoreCase, "i", false, "игнорировать регистр")
	flag.BoolVar(&opts.Invert, "v", false, "вывести строки, не содержащие шаблон")
	flag.BoolVar(&opts.Fixed, "F", false, "воспринимать шаблон как фиксированную строку")
	flag.BoolVar(&opts.LineNum, "n", false, "вывести номер строки перед каждой найденной строкой")

	flag.Parse()

	// проверяем, что указан шаблон
	if flag.NArg() < 1 {
		return nil, fmt.Errorf("должен быть указан шаблон")
	}

	opts.Pattern = flag.Arg(0)
	files := flag.Args()[1:]

	// если не указаны файлы, то используем stdin
	if len(files) == 0 {
		files = []string{"-"}
	}

	return &models.GrepConfig{
		Options: opts,
		Files:   files,
	}, nil
}
