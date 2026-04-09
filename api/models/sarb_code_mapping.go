package models

import "time"

// SARBCodeMapping maps an IFRS 17 line item to a PA regulatory return code.
type SARBCodeMapping struct {
	ID                     int       `json:"id" gorm:"primaryKey;autoIncrement"`
	LineItem               string    `json:"line_item"`               // e.g. "Insurance Contract Liabilities"
	Description            string    `json:"description"`
	IFRS17BalanceSheetItem string    `json:"ifrs17_balance_sheet_item"`
	RegulatoryReturn       string    `json:"regulatory_return"`       // e.g. "RC&AP" or "J200"
	ReturnCode             string    `json:"return_code"`             // e.g. "RC01", "J200-L01"
	DebitCreditIndicator   string    `json:"debit_credit_indicator"`  // "D" or "C"
	Notes                  string    `json:"notes"`
	CreatedBy              string    `json:"created_by"`
	CreatedAt              time.Time `json:"created_at"`
	UpdatedAt              time.Time `json:"updated_at"`
}

// SARBReportRow is a generated regulatory return row combining mapping with actual values.
type SARBReportRow struct {
	ReturnCode           string  `json:"return_code"`
	RegulatoryReturn     string  `json:"regulatory_return"`
	LineItem             string  `json:"line_item"`
	Description          string  `json:"description"`
	DebitCreditIndicator string  `json:"debit_credit_indicator"`
	Amount               float64 `json:"amount"`
	Notes                string  `json:"notes"`
}

type UpsertSARBMappingRequest struct {
	ID                     int    `json:"id"`
	LineItem               string `json:"line_item"`
	Description            string `json:"description"`
	IFRS17BalanceSheetItem string `json:"ifrs17_balance_sheet_item"`
	RegulatoryReturn       string `json:"regulatory_return"`
	ReturnCode             string `json:"return_code"`
	DebitCreditIndicator   string `json:"debit_credit_indicator"`
	Notes                  string `json:"notes"`
}
