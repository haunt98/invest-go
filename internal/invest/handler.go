package invest

import (
	"context"
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
)

type Handler interface {
	List(ctx context.Context) error
}

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{
		service: service,
	}
}

func (h *handler) List(ctx context.Context) error {
	investments, err := h.service.List(ctx)
	if err != nil {
		return fmt.Errorf("service failed to list: %w", err)
	}

	// https://github.com/jedib0t/go-pretty/tree/main/table
	tableWriter := table.NewWriter()
	tableWriter.SetOutputMirror(os.Stdout)
	tableWriter.AppendHeader(table.Row{
		"ID",
		"Amount",
		"Date",
		"Source",
	})

	totalAmount := 0
	for _, investment := range investments {
		tableWriter.AppendRow(table.Row{
			investment.ID,
			investment.Amount,
			investment.Date,
			investment.Source,
		})
		totalAmount += int(investment.Amount)
	}
	tableWriter.AppendSeparator()
	tableWriter.AppendFooter(table.Row{
		"",
		totalAmount,
		"",
		"",
	})
	tableWriter.Render()

	return nil
}
