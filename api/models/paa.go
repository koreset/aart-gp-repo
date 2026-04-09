package models

type PAARun struct {
	RunName       string `json:"run_name"`
	PortfolioName string `json:"portfolio_name"`
	RunID         int    `json:"run_id"`
}

type TableYearVersion struct {
	Year    int      `json:"year"`
	Version []string `json:"version"`
}
