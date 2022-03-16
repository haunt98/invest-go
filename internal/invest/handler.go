package invest

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/make-go-great/ioe-go"
	"github.com/spf13/cast"
)

type Handler interface {
	List(ctx context.Context) error
	Add(ctx context.Context, investment Investment, interactive bool) error
	Remove(ctx context.Context, id string, interactive bool) error
	Export(ctx context.Context, filename string) error
	Import(ctx context.Context, filename string) error
}

type handler struct {
	service Service
}

func NewHandler(service Service) Handler {
	return &handler{
		service: service,
	}
}

func (h *handler) List(ctx context.Context) error {
	investments, err := h.service.List(ctx)
	if err != nil {
		return err
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

func (h *handler) Add(ctx context.Context, investment Investment, interactive bool) error {
	if interactive {
		fmt.Printf("Input amount: ")
		investment.Amount = cast.ToInt64(ioe.ReadInput())

		fmt.Printf("Input date: ")
		investment.Date = ioe.ReadInput()

		fmt.Printf("Input source: ")
		investment.Source = ioe.ReadInput()
	}

	return h.service.Add(ctx, investment)
}

func (h *handler) Remove(ctx context.Context, id string, interactive bool) error {
	if interactive {
		fmt.Printf("Input ID: ")
		id = ioe.ReadInput()
	}

	return h.service.Remove(ctx, id)
}

func (h *handler) Export(ctx context.Context, filename string) error {
	investments, err := h.service.List(ctx)
	if err != nil {
		return fmt.Errorf("service failed to list: %w", err)
	}

	data := WrapInvestments{
		Investments: investments,
	}

	bytes, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return fmt.Errorf("json failed to marshal indent: %w", err)
	}

	if err := os.WriteFile(filename, bytes, 0o755); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

func (h *handler) Import(ctx context.Context, filename string) error {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	data := WrapInvestments{}

	if err := json.Unmarshal(bytes, &data); err != nil {
		return fmt.Errorf("json failed to unmarshal: %w", err)
	}

	for _, investment := range data.Investments {
		if err := h.service.Add(ctx, investment); err != nil {
			return fmt.Errorf("service failed to add: %w", err)
		}
	}

	return nil
}
