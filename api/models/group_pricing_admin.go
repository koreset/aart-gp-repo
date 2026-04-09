package models

type MemberBeneficiary struct {
	ID             int
	BeneficiaryID  int     `json:"beneficiary_id" gorm:"column:beneficiary_id"`
	MemberID       int     `json:"member_id" gorm:"column:member_id"`
	Relationship   string  `json:"relationship"`
	Name           string  `json:"name"`
	Percentage     float64 `json:"percentage"`
	IDNumber       string  `json:"id_number"`
	PassportNumber string  `json:"passport_number"`
	ContactNumber  string  `json:"contact_number"`
	Email          string  `json:"email"`
}
