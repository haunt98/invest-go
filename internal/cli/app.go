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
	name  = "invest"
	usage = "tracking invest money"

	commandList   = "list"
	commandAdd    = "add"
	commandRemove = "remove"

	flagVerbose = "verbose"
	flagID      = "id"
	flagAmount  = "amount"
	flagDate    = "date"
	flagSource  = "source"

	usageList    = "list all investments"
	usageAdd     = "add investment"
	usageRemove  = "remove investment"
	usageVerbose = "debug"
	usageID      = "id of investment"
	usageAmount  = "amount of investment"
	usageDate    = "date of investment, example 2022-12-31"
	usageSource  = "source of investment"
)

var aliasesVerbose = []string{"v"}

type App struct {
	cliApp *cli.App
}

func NewApp(db *sql.DB, shouldInitDatabase bool, location *time.Location) (*App, error) {
	investRepo, err := invest.NewRepository(context.Background(), db, shouldInitDatabase)
	if err != nil {
		return nil, fmt.Errorf("failed to new repository: %w", err)
	}
	investService := invest.NewService(investRepo, location)
	investHandler := invest.NewHandler(investService)

	a := &action{
		investHandler: investHandler,
	}

	cliApp := &cli.App{
		Name:   name,
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
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    flagVerbose,
						Aliases: aliasesVerbose,
						Usage:   usageVerbose,
					},
					&cli.Int64Flag{
						Name:     flagAmount,
						Usage:    usageAmount,
						Required: true,
					},
					&cli.StringFlag{
						Name:     flagDate,
						Usage:    usageDate,
						Required: true,
					},
					&cli.StringFlag{
						Name:     flagSource,
						Usage:    usageSource,
						Required: true,
					},
				},
			},
			{
				Name:   commandRemove,
				Usage:  usageRemove,
				Action: a.RunRemove,
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    flagVerbose,
						Aliases: aliasesVerbose,
						Usage:   usageVerbose,
					},
					&cli.StringFlag{
						Name:     flagID,
						Usage:    usageID,
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
		color.PrintAppError(name, err.Error())
	}
}
