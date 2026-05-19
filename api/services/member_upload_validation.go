package services

import "fmt"

// MemberUploadBlockingError describes a single row in a member upload that is
// missing a hard-required field (gender, date_of_birth, annual_salary). These
// fields drive pricing and valuation maths, so a row missing any of them must
// not be persisted.
type MemberUploadBlockingError struct {
	Row            int    `json:"row"`
	Field          string `json:"field"`
	Message        string `json:"message"`
	MemberIdNumber string `json:"member_id_number,omitempty"`
	MemberName     string `json:"member_name,omitempty"`
}

// MemberUploadBlockingErrors is an error type returned by member-upload services
// when one or more rows are missing required fields. The controller layer should
// type-assert via errors.As and respond with HTTP 400 plus the structured list
// so the frontend can render a per-row report and a downloadable error CSV.
type MemberUploadBlockingErrors struct {
	Errors []MemberUploadBlockingError
}

func (e *MemberUploadBlockingErrors) Error() string {
	if e == nil || len(e.Errors) == 0 {
		return "member upload validation failed"
	}
	return fmt.Sprintf(
		"member upload blocked: %d row(s) missing required fields (gender, date_of_birth, annual_salary)",
		len(e.Errors),
	)
}
