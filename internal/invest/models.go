package invest

type Investment struct {
	ID     string `json:"id"`
	Amount int64  `json:"amount"`
	Date   string `json:"date"`
	Source string `json:"source"`
}

type WrapInvestments struct {
	Investments []Investment `json:"investments"`
}
