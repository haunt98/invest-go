package cli

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/haunt98/invest-go/internal/invest"
	"github.com/make-go-great/color-go"
	"github.com/urfave/cli/v2"
)

const (
	Name  = "invest"
	usage = "tracking invest money"

	commandList   = "list"
	commandAdd    = "add"
	commandRemove = "remove"
	commandExport = "export"
	commandImport = "import"

	flagFilename = "filename"

	usageList     = "list all investments"
	usageAdd      = "add investment"
	usageRemove   = "remove investment"
	usageExport   = "export data"
	usageImport   = "import data"
	usageFilename = "filename to export/import"
)

var aliasInteractive = []string{"i"}

type App struct {
	cliApp *cli.App
}

func NewApp(db *sql.DB, location *time.Location) (*App, error) {
	investRepo, err := invest.NewRepository(context.Background(), db)
	if err != nil {
		return nil, fmt.Errorf("failed to new repository: %w", err)
	}
	investService := invest.NewService(investRepo, location)
	investHandler := invest.NewHandler(investService)

	a := &action{
		investHandler: investHandler,
	}

	cliApp := &cli.App{
		Name:   Name,
		Usage:  usage,
		Action: a.RunHelp,
		Commands: []*cli.Command{
			{
				Name:   commandList,
				Usage:  usageList,
				Action: a.RunList,
			},
			{
				Name:   commandAdd,
				Usage:  usageAdd,
				Action: a.RunAdd,
			},
			{
				Name:   commandRemove,
				Usage:  usageRemove,
				Action: a.RunRemove,
			},
			{
				Name:   commandExport,
				Usage:  usageExport,
				Action: a.RunExport,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     flagFilename,
						Usage:    usageFilename,
						Required: true,
					},
				},
			},
			{
				Name:   commandImport,
				Usage:  usageImport,
				Action: a.RunImport,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     flagFilename,
						Usage:    usageFilename,
						Required: true,
					},
				},
			},
		},
	}

	return &App{
		cliApp: cliApp,
	}, nil
}

func (a *App) Run() {
	if err := a.cliApp.Run(os.Args); err != nil {
		color.PrintAppError(Name, err.Error())
	}
}
