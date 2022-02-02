package invest

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Service interface {
	List(ctx context.Context) ([]Investment, error)
	Add(ctx context.Context, investment Investment) error
	Remove(ctx context.Context, id string) error
}

type service struct {
	repo     Repository
	location *time.Location
}

func NewService(
	repo Repository,
	location *time.Location,
) Service {
	return &service{
		repo:     repo,
		location: location,
	}
}

func (s *service) List(ctx context.Context) ([]Investment, error) {
	investments, err := s.repo.GetInvestments(ctx)
	if err != nil {
		return nil, fmt.Errorf("repository failed to get investments: %w", err)
	}

	// Simple date format to output
	for i := range investments {
		investments[i].Date, err = dateToOutput(investments[i].Date, s.location)
		if err != nil {
			return nil, err
		}
	}

	return investments, nil
}

func (s *service) Add(ctx context.Context, investment Investment) error {
	if investment.ID == "" {
		investment.ID = uuid.NewString()
	}

	if !validateInvestment(investment) {
		return fmt.Errorf("investment %+v not valid", investment)
	}

	// Better date format from input
	var err error
	investment.Date, err = dateFromInput(investment.Date, s.location)
	if err != nil {
		return err
	}

	if err := s.repo.CreateInvestment(ctx, investment); err != nil {
		return fmt.Errorf("repository failed to create investment: %w", err)
	}

	return nil
}

func (s *service) Remove(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("id not valid")
	}

	if err := s.repo.DeleteInvestment(ctx, id); err != nil {
		return fmt.Errorf("repository failed to delete investment: %w", err)
	}

	return nil
}

func validateInvestment(investment Investment) bool {
	if investment.ID == "" {
		return false
	}

	if investment.Amount == 0 {
		return false
	}

	if investment.Date == "" {
		return false
	}

	if investment.Source == "" {
		return false
	}

	return true
}
