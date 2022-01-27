package cli

import (
	"github.com/haunt98/invest-go/internal/invest"
	"github.com/urfave/cli/v2"
)

type action struct {
	// flags struct {
	// 	amount string
	// 	date   string
	// 	source string
	// }

	investHandler invest.Handler
}

func (a *action) RunHelp(c *cli.Context) error {
	return cli.ShowAppHelp(c)
}

func (a *action) RunList(c *cli.Context) error {
	return a.investHandler.List(c.Context)
}
