package cli

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/haunt98/invest-go/internal/invest"
	"github.com/make-go-great/color-go"
	"github.com/urfave/cli/v2"
)

const (
	appName  = "invest"
	appUsage = "tracking invest money"

	// commands
	listCommand = "list"

	// command usage
	listUsage = "list all investments"
)

type App struct {
	cliApp *cli.App
}

func NewApp(db *sql.DB) (*App, error) {
	investRepo, err := invest.NewRepository(context.Background(), db)
	if err != nil {
		return nil, fmt.Errorf("failed to new repository: %w", err)
	}
	investService := invest.NewService(investRepo)
	investHandler := invest.NewHandler(investService)

	a := &action{
		investHandler: investHandler,
	}

	cliApp := &cli.App{
		Name:  appName,
		Usage: appUsage,
		Commands: []*cli.Command{
			{
				Name:   listCommand,
				Usage:  listUsage,
				Action: a.RunList,
			},
		},
		Action: a.RunHelp,
	}

	return &App{
		cliApp: cliApp,
	}, nil
}

func (a *App) Run() {
	if err := a.cliApp.Run(os.Args); err != nil {
		color.PrintAppError(appName, err.Error())
	}
}
