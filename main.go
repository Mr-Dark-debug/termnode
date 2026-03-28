package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/Mr-Dark-debug/termnode/internal/app"
	"github.com/Mr-Dark-debug/termnode/internal/db"
)

var version = "dev"

func main() {
	dbPath := flag.String("db", defaultDBPath(), "path to SQLite database")
	port := flag.String("port", ":8080", "HTTP webhook listen address")
	debug := flag.Bool("debug", false, "enable debug logging")
	showVersion := flag.Bool("version", false, "print version and exit")
	flag.Parse()

	if *showVersion {
		fmt.Printf("termnode %s\n", version)
		os.Exit(0)
	}

	if *debug {
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to open log: %v\n", err)
			os.Exit(1)
		}
		defer f.Close()
	}

	database, err := db.Open(*dbPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open database: %v\n", err)
		os.Exit(1)
	}
	defer database.Close()

	repo := db.NewRepository(database)
	appModel := app.New(repo, *port, version)

	p := tea.NewProgram(appModel, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error running program: %v\n", err)
		os.Exit(1)
	}
}

func defaultDBPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "/data/data/com.termux/files/home"
	}
	return filepath.Join(home, ".termnode", "termnode.db")
}
