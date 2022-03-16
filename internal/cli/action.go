package cli

import (
	"github.com/haunt98/invest-go/internal/invest"
	"github.com/urfave/cli/v2"
)

type action struct {
	flags struct {
		// Invest
		id     string
		amount int64
		date   string
		source string

		// Export, import
		filename string

		// Interactive mode
		interactive bool
	}

	investHandler invest.Handler
}

func (a *action) RunHelp(c *cli.Context) error {
	return cli.ShowAppHelp(c)
}

func (a *action) RunList(c *cli.Context) error {
	return a.investHandler.List(c.Context)
}

func (a *action) RunAdd(c *cli.Context) error {
	a.getFlags(c)

	return a.investHandler.Add(c.Context, invest.Investment{
		Amount: a.flags.amount,
		Date:   a.flags.date,
		Source: a.flags.source,
	}, a.flags.interactive)
}

func (a *action) RunRemove(c *cli.Context) error {
	a.getFlags(c)

	return a.investHandler.Remove(c.Context, a.flags.id, a.flags.interactive)
}

func (a *action) RunExport(c *cli.Context) error {
	a.getFlags(c)

	return a.investHandler.Export(c.Context, a.flags.filename)
}

func (a *action) RunImport(c *cli.Context) error {
	a.getFlags(c)

	return a.investHandler.Import(c.Context, a.flags.filename)
}

func (a *action) getFlags(c *cli.Context) {
	a.flags.id = c.String(flagID)
	a.flags.amount = c.Int64(flagAmount)
	a.flags.date = c.String(flagDate)
	a.flags.source = c.String(flagSource)
	a.flags.filename = c.String(flagFilename)
	a.flags.interactive = c.Bool(flagInteractive)
}
