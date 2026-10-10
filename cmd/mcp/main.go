package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/sqot0/packsmith/backend/cmd"
)

func main() {
	projectPath := flag.String("project", "", "Path to modpack project directory")
	flag.Parse()

	// If positional arg is provided, use it as project path
	if *projectPath == "" && flag.NArg() > 0 {
		*projectPath = flag.Arg(0)
	}

	app := cmd.NewApp()
	if *projectPath != "" {
		if _, err := app.OpenProject(*projectPath); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to open initial project at %s: %v\n", *projectPath, err)
		}
	}

	if err := app.ServeMCP(); err != nil {
		fmt.Fprintf(os.Stderr, "MCP server error: %v\n", err)
		os.Exit(1)
	}
}
