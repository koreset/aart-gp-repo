package models

import "time"

// Beneficiary represents a member's beneficiary record.
// Table: member_beneficiaries
type Beneficiary struct {
    ID                  int         `json:"id" gorm:"primaryKey;column:id"`
    FullName            string      `json:"full_name" gorm:"column:full_name;size:255;not null"`
    Relationship        string      `json:"relationship" gorm:"column:relationship;size:100;not null"`
    IDType              string      `json:"id_type" gorm:"column:id_type;size:50;not null"`
    IDNumber            string      `json:"id_number" gorm:"column:id_number;size:100;not null"`
    Gender              string      `json:"gender" gorm:"column:gender;size:20;not null"`
    DateOfBirth         *time.Time  `json:"date_of_birth" gorm:"column:date_of_birth"`
    ContactNumber       string      `json:"contact_number,omitempty" gorm:"column:contact_number;size:100"`
    Email               string      `json:"email,omitempty" gorm:"column:email;size:255"`
    Address             string      `json:"address,omitempty" gorm:"column:address;size:500"`
    AllocationPercentage float64    `json:"allocation_percentage" gorm:"column:allocation_percentage;not null"`
    BenefitTypes        StringArray `json:"benefit_types" gorm:"column:benefit_types;type:json"`
    GuardianName        string      `json:"guardian_name,omitempty" gorm:"column:guardian_name;size:255"`
    GuardianRelationship string     `json:"guardian_relationship,omitempty" gorm:"column:guardian_relationship;size:100"`
    GuardianIDNumber    string      `json:"guardian_id_number,omitempty" gorm:"column:guardian_id_number;size:100"`
    GuardianContact     string      `json:"guardian_contact,omitempty" gorm:"column:guardian_contact;size:100"`
    BankName            string      `json:"bank_name,omitempty" gorm:"column:bank_name;size:255"`
    BranchCode          string      `json:"branch_code,omitempty" gorm:"column:branch_code;size:50"`
    AccountNumber       string      `json:"account_number,omitempty" gorm:"column:account_number;size:100"`
    AccountType         string      `json:"account_type,omitempty" gorm:"column:account_type;size:50"`
    Status              string      `json:"status,omitempty" gorm:"column:status;size:50"`
    MemberID            int         `json:"memberId" gorm:"column:member_id;index;not null"`
    CreatedAt           time.Time   `json:"created_at" gorm:"column:created_at;autoCreateTime"`
    UpdatedAt           time.Time   `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
}

func (Beneficiary) TableName() string { return "member_beneficiaries" }
