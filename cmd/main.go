package main

import (
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/urfave/cli/v2"
	"github.com/woojiahao/baleen/internal/baleen"
)

// TODO: Support general migrations from Trello to Notion
func main() {
	var boardName, envPath, configPath, savePath string
	var toSave bool

	app := &cli.App{
		Name:  "baleen",
		Usage: "migrate your Trello thoughts board to Notion",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "board",
				Aliases:     []string{"b"},
				Value:       "Programming Bucket",
				Usage:       "specify name of Trello board",
				Destination: &boardName,
			},
			&cli.StringFlag{
				Name:        "env",
				Aliases:     []string{"e"},
				Value:       ".env",
				Usage:       "specify the environment file hodling the API keys",
				Destination: &envPath,
			},
			&cli.StringFlag{
				Name:        "config",
				Aliases:     []string{"c"},
				Value:       "configs/conf.json",
				Usage:       "specify the configuration JSON for the migration",
				Destination: &configPath,
			},
		},
		Commands: []*cli.Command{
			{
				Name: "migrate",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:        "save",
						Aliases:     []string{"s"},
						Value:       true,
						Usage:       "specify whether to save files during migration (used in \"baleen migrate\")",
						Destination: &toSave,
					},
				},
				Usage: "exports a Trello board into the integrated Notion page (full flow from saving exports to importing to Notion)",
				Action: func(c *cli.Context) error {
					baleen.Migrate(boardName, configPath, envPath, toSave)
					return nil
				},
			},
			{
				Name:  "import",
				Usage: "imports saved cards into Notion (cards saved from \"baleen migrate\" or \"baleen export\")",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "savePath",
						Aliases:     []string{"sp"},
						Usage:       "specify the path of a save file",
						Destination: &savePath,
					},
				},
				Action: func(c *cli.Context) error {
					if savePath == "" {
						return fmt.Errorf("save path not specified")
					}
					baleen.Import(savePath, configPath, envPath)
					return nil
				},
			},
			{
				Name:  "export",
				Usage: "exports a Trello board and creates a save file (to import, use \"baleen import <save path>\"",
				Action: func(c *cli.Context) error {
					baleen.ExportAndSave(boardName, envPath)
					return nil
				},
			},
			{
				Name:  "archive",
				Usage: "archives all cards in lists in Trello",
				Action: func(c *cli.Context) error {
					baleen.ClearBoard(boardName, envPath)
					return nil
				},
			},
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	err := app.Run(os.Args)
	if err != nil {
		log.Fatalf("Error running cli: %v\n", err)
	}
}
