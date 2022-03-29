package invest

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"
)

type RepoSuite struct {
	suite.Suite

	db *sql.DB

	dbMock        sqlmock.Sqlmock
	preparedStmts map[string]*sqlmock.ExpectedPrepare
	columns       []string

	repo Repository
}

func (s *RepoSuite) SetupTest() {
	db, dbMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	s.NoError(err)
	s.db = db
	s.dbMock = dbMock
	s.columns = []string{"id", "amount", "date", "source"}

	dbMock.ExpectExec(stmtInitInvestments).
		WillReturnResult(sqlmock.NewResult(1, 1))

	s.preparedStmts = make(map[string]*sqlmock.ExpectedPrepare)
	s.preparedStmts[preparedCreateInvestment] = s.dbMock.ExpectPrepare(stmtCreateInvestment)
	s.preparedStmts[preparedDeleteInvestment] = s.dbMock.ExpectPrepare(stmtDeleteInvestment)

	s.repo, err = NewRepository(context.Background(), s.db)
	s.NoError(err)
}

func TestRepoSuite(t *testing.T) {
	suite.Run(t, new(RepoSuite))
}

func (s *RepoSuite) TestInitDatabaseFailed() {
	db, dbMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	s.NoError(err)

	dbMock.ExpectExec(stmtInitInvestments).
		WillReturnError(errors.New("failed"))

	_, err = NewRepository(context.Background(), db)
	s.Error(err)
}

func (s *RepoSuite) TestGetInvestmentsSuccess() {
	investments := []Investment{
		{
			ID:     "1",
			Amount: 1000,
			Date:   "date_1",
			Source: "source_1",
		},
		{
			ID:     "2",
			Amount: 2000,
			Date:   "date_2",
			Source: "source_2",
		},
	}

	rows := sqlmock.NewRows(s.columns)
	for _, investment := range investments {
		rows.AddRow(investment.ID, investment.Amount, investment.Date, investment.Source)
	}

	s.dbMock.ExpectQuery(stmtGetInvestments).
		WillReturnRows(rows)

	gotResult, gotErr := s.repo.GetInvestments(context.Background())
	s.NoError(gotErr)
	s.Equal(investments, gotResult)
}

func (s *RepoSuite) TestGetInvestmentsFailedDatabase() {
	investments := []Investment{
		{
			ID:     "1",
			Amount: 1000,
			Date:   "date_1",
			Source: "source_1",
		},
		{
			ID:     "2",
			Amount: 2000,
			Date:   "date_2",
			Source: "source_2",
		},
	}

	rows := sqlmock.NewRows(s.columns)
	for _, investment := range investments {
		rows.AddRow(investment.ID, investment.Amount, investment.Date, investment.Source)
	}

	s.dbMock.ExpectQuery(stmtGetInvestments).
		WillReturnError(errors.New("failed"))

	_, gotErr := s.repo.GetInvestments(context.Background())
	s.Error(gotErr)
}

func (s *RepoSuite) TestCreateInvestmentSuccess() {
	investment := Investment{
		ID:     "1",
		Amount: 10,
		Date:   "date",
		Source: "source",
	}

	s.preparedStmts[preparedCreateInvestment].ExpectExec().
		WithArgs(
			investment.ID,
			investment.Amount,
			investment.Date,
			investment.Source,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	gotErr := s.repo.CreateInvestment(context.Background(), investment)
	s.NoError(gotErr)
}

func (s *RepoSuite) TestCreateInvestmentFailedDatabase() {
	investment := Investment{
		ID:     "1",
		Amount: 10,
		Date:   "date",
		Source: "source",
	}

	s.preparedStmts[preparedCreateInvestment].ExpectExec().
		WithArgs(
			investment.ID,
			investment.Amount,
			investment.Date,
			investment.Source,
		).
		WillReturnError(errors.New("failed"))

	gotErr := s.repo.CreateInvestment(context.Background(), investment)
	s.Error(gotErr)
}

func (s *RepoSuite) TestDeleteInvestmentSuccess() {
	id := "1"

	s.preparedStmts[preparedDeleteInvestment].ExpectExec().
		WithArgs(
			id,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	gotErr := s.repo.DeleteInvestment(context.Background(), id)
	s.NoError(gotErr)
}

func (s *RepoSuite) TestDeleteInvestmentFailedDatabase() {
	id := "1"

	s.preparedStmts[preparedDeleteInvestment].ExpectExec().
		WithArgs(
			id,
		).
		WillReturnError(errors.New("failed"))

	gotErr := s.repo.DeleteInvestment(context.Background(), id)
	s.Error(gotErr)
}
