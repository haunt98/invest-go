package cli

import (
	"github.com/haunt98/invest-go/internal/invest"
	"github.com/urfave/cli/v2"
)

type action struct {
	flags struct {
		// Export, import
		filename string
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
	return a.investHandler.Add(c.Context)
}

func (a *action) RunRemove(c *cli.Context) error {
	return a.investHandler.Remove(c.Context)
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
	a.flags.filename = c.String(flagFilename)
}
