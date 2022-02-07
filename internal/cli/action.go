package cli

import (
	"errors"
	"fmt"

	"github.com/haunt98/invest-go/internal/invest"
	"github.com/make-go-great/ioe-go"
	"github.com/spf13/cast"
	"github.com/urfave/cli/v2"
)

type action struct {
	flags struct {
		// Debug
		verbose bool

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

	if !a.flags.interactive {
		if a.flags.amount == 0 {
			return errors.New("empty amount")
		}

		if a.flags.date == "" {
			return errors.New("empty date")
		}

		if a.flags.source == "" {
			return errors.New("empty source")
		}
	} else {
		fmt.Printf("Input amount (%s):\n", usageAmount)
		a.flags.amount = cast.ToInt64(ioe.ReadInput())

		fmt.Printf("Input date (%s):\n", usageDate)
		a.flags.date = ioe.ReadInput()

		fmt.Printf("Input source (%s):\n", usageSource)
		a.flags.source = ioe.ReadInput()
	}

	return a.investHandler.Add(c.Context, invest.Investment{
		Amount: a.flags.amount,
		Date:   a.flags.date,
		Source: a.flags.source,
	})
}

func (a *action) RunRemove(c *cli.Context) error {
	a.getFlags(c)

	if !a.flags.interactive {
		if a.flags.id == "" {
			return errors.New("empty id")
		}
	} else {
		fmt.Printf("Input id (%s):\n", usageID)
		a.flags.id = ioe.ReadInput()
	}

	return a.investHandler.Remove(c.Context, a.flags.id)
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
	a.flags.verbose = c.Bool(flagVerbose)
	a.flags.id = c.String(flagID)
	a.flags.amount = c.Int64(flagAmount)
	a.flags.date = c.String(flagDate)
	a.flags.source = c.String(flagSource)
	a.flags.filename = c.String(flagFilename)
	a.flags.interactive = c.Bool(flagInteractive)
}
