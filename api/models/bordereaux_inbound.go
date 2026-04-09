package models

import "time"

// EmployerSubmission represents an inbound member data submission from an employer.
// Status lifecycle: pending_receipt → received → under_review → (queries_raised →) accepted | rejected
type EmployerSubmission struct {
	ID                      int        `gorm:"primaryKey;autoIncrement" json:"id"`
	SchemeID                int        `json:"scheme_id" gorm:"not null;index"`
	SchemeName              string     `json:"scheme_name" gorm:"default:''"`
	Month                   int        `json:"month" gorm:"not null"`
	Year                    int        `json:"year" gorm:"not null"`
	DueDate                 string     `json:"due_date" gorm:"default:''"`
	ReceivedDate            *time.Time `json:"received_date"`
	SubmittedBy             string     `json:"submitted_by" gorm:"default:''"`
	// Status: pending_receipt | received | under_review | queries_raised | accepted | rejected
	Status                  string     `json:"status" gorm:"default:'pending_receipt'"`
	FileName                string     `json:"file_name" gorm:"default:''"`
	FilePath                string     `json:"file_path" gorm:"default:''"`
	FileSize                int64      `json:"file_size" gorm:"default:0"`
	RecordCount             int        `json:"record_count" gorm:"default:0"`
	ValidCount              int        `json:"valid_count" gorm:"default:0"`
	InvalidCount            int        `json:"invalid_count" gorm:"default:0"`
	Notes                   string     `json:"notes" gorm:"type:text"`
	QueryNotes              string     `json:"query_notes" gorm:"type:text"`
	LinkedPremiumScheduleID *int       `json:"linked_premium_schedule_id"`
	ReviewedBy              string     `json:"reviewed_by" gorm:"default:''"`
	ReviewedAt              *time.Time `json:"reviewed_at"`
	AcceptedBy              string     `json:"accepted_by" gorm:"default:''"`
	AcceptedAt              *time.Time `json:"accepted_at"`
	RejectedBy              string     `json:"rejected_by" gorm:"default:''"`
	RejectedAt              *time.Time `json:"rejected_at"`
	RejectionReason         string     `json:"rejection_reason" gorm:"default:''"`
	// Retro submission fields
	IsRetro             bool   `json:"is_retro" gorm:"default:false"`
	RetroEffectiveDate  string `json:"retro_effective_date" gorm:"default:''"`
	// Member register sync (legacy single-pass)
	MemberSyncedAt      *time.Time `json:"member_synced_at"`
	MemberSyncSummary   string     `json:"member_sync_summary" gorm:"type:text"`
	// Per-category register sync tracking
	ExitsSyncedAt         *time.Time `json:"exits_synced_at"`
	ExitsSyncSummary      string     `json:"exits_sync_summary" gorm:"type:text"`
	AmendmentsSyncedAt    *time.Time `json:"amendments_synced_at"`
	AmendmentsSyncSummary string     `json:"amendments_sync_summary" gorm:"type:text"`
	NewJoinersSyncedAt    *time.Time `json:"new_joiners_synced_at"`
	NewJoinersSyncSummary string     `json:"new_joiners_sync_summary" gorm:"type:text"`
	CreatedBy               string     `json:"created_by" gorm:"default:''"`
	CreatedAt               time.Time  `json:"created_at"`
	UpdatedAt               time.Time  `json:"updated_at"`
	Records                 []EmployerSubmissionRecord `json:"records" gorm:"foreignKey:SubmissionID"`
}

// EmployerSubmissionRecord is a single member row parsed from an employer submission file.
type EmployerSubmissionRecord struct {
	ID               int     `gorm:"primaryKey;autoIncrement" json:"id"`
	SubmissionID     int     `json:"submission_id" gorm:"not null;index"`
	RowNumber        int     `json:"row_number" gorm:"default:0"`
	MemberName       string  `json:"member_name" gorm:"default:''"`
	EmployeeNumber   string  `json:"employee_number" gorm:"default:''"`
	IDNumber         string  `json:"id_number" gorm:"default:''"`
	IDType           string  `json:"id_type" gorm:"default:''"` // "rsa_id" | "passport" — auto-detected on upload
	DateOfBirth      string  `json:"date_of_birth" gorm:"default:''"`
	Salary           float64 `json:"salary" gorm:"default:0"`
	BenefitCode      string  `json:"benefit_code" gorm:"default:''"`
	PremiumAmount    float64 `json:"premium_amount" gorm:"default:0"`
	EntryDate        string  `json:"entry_date" gorm:"default:''"`
	ExitDate         string  `json:"exit_date" gorm:"default:''"`
	Gender           string  `json:"gender" gorm:"default:''"`
	// ValidationStatus: valid | id_invalid | amount_invalid | missing_data | excluded
	ValidationStatus string  `json:"validation_status" gorm:"default:'valid'"`
	ExclusionReason  string  `json:"exclusion_reason" gorm:"default:''"`
	RawData          string  `json:"raw_data" gorm:"type:text"`
}

// CreateSubmissionRequest is the request body for creating a new employer submission record.
type CreateSubmissionRequest struct {
	SchemeID           int    `json:"scheme_id" binding:"required"`
	Month              int    `json:"month" binding:"required"`
	Year               int    `json:"year" binding:"required"`
	DueDate            string `json:"due_date"`
	SubmittedBy        string `json:"submitted_by"`
	Notes              string `json:"notes"`
	IsRetro            bool   `json:"is_retro"`
	RetroEffectiveDate string `json:"retro_effective_date"`
}

// RaiseQueryRequest is the request body for raising a query on a submission.
type RaiseQueryRequest struct {
	QueryNotes string `json:"query_notes" binding:"required"`
}

// RejectSubmissionRequest is the request body for rejecting a submission.
type RejectSubmissionRequest struct {
	Reason string `json:"reason" binding:"required"`
}

// SubmissionDeltaRecord stores the computed change classification for a single member
// between the current accepted submission and the prior period's accepted submission.
// ChangeType values: new | amendment | ceased | continuing
type SubmissionDeltaRecord struct {
	ID                int    `gorm:"primaryKey;autoIncrement" json:"id"`
	SubmissionID      int    `json:"submission_id" gorm:"not null;index"`
	PriorSubmissionID int    `json:"prior_submission_id" gorm:"default:0"`
	MemberKey         string `json:"member_key" gorm:"not null"`
	MemberName        string `json:"member_name" gorm:"default:''"`
	ChangeType        string `json:"change_type" gorm:"not null"` // new | amendment | ceased | continuing
	ChangedFields     string `json:"changed_fields" gorm:"type:text"`
	RowNumber         int    `json:"row_number" gorm:"default:0"`
}

// MemberSyncResult is returned by SyncSubmissionToMemberRegister and stored as JSON on the submission.
type MemberSyncResult struct {
	Added       int      `json:"added"`
	Updated     int      `json:"updated"`
	Deactivated int      `json:"deactivated"`
	Skipped     int      `json:"skipped"`
	Errors      []string `json:"errors"`
}

// DeltaSummary is the aggregate count returned by ComputeSubmissionDelta.
type DeltaSummary struct {
	SubmissionID      int `json:"submission_id"`
	PriorSubmissionID int `json:"prior_submission_id"`
	New               int `json:"new"`
	Amendment         int `json:"amendment"`
	Ceased            int `json:"ceased"`
	Continuing        int `json:"continuing"`
	Total             int `json:"total"`
}

// RegisterDiffMember is a single member entry in a register diff result.
type RegisterDiffMember struct {
	MemberKey      string              `json:"member_key"`
	MemberName     string              `json:"member_name"`
	EmployeeNumber string              `json:"employee_number"`
	IDNumber       string              `json:"id_number"`
	RowNumber      int                 `json:"row_number"`
	ExitDate       string              `json:"exit_date,omitempty"`
	ChangedFields  map[string][2]string `json:"changed_fields,omitempty"`
}

// RegisterDiffResult is the summary returned by ComputeRegisterDiff.
type RegisterDiffResult struct {
	SubmissionID int                  `json:"submission_id"`
	NewJoiners   []RegisterDiffMember `json:"new_joiners"`
	Exits        []RegisterDiffMember `json:"exits"`
	Amendments   []RegisterDiffMember `json:"amendments"`
	Continuing   int                  `json:"continuing"`
	// Snapshot provenance — populated when result is served from a persisted snapshot.
	IsSnapshot bool       `json:"is_snapshot"`
	SnapshotAt *time.Time `json:"snapshot_at,omitempty"`
	SnapshotBy string     `json:"snapshot_by,omitempty"`
}

// RegisterDiffSnapshot persists the full register diff at the moment a submission is accepted.
// This ensures the historical diff is always available even after subsequent register updates.
type RegisterDiffSnapshot struct {
	ID              int       `gorm:"primaryKey;autoIncrement" json:"id"`
	SubmissionID    int       `json:"submission_id" gorm:"not null;uniqueIndex"`
	SnapshotAt      time.Time `json:"snapshot_at"`
	SnapshotBy      string    `json:"snapshot_by" gorm:"default:''"`
	NewJoinersCount int       `json:"new_joiners_count" gorm:"default:0"`
	ExitsCount      int       `json:"exits_count" gorm:"default:0"`
	AmendmentsCount int       `json:"amendments_count" gorm:"default:0"`
	ContinuingCount int       `json:"continuing_count" gorm:"default:0"`
	// Full RegisterDiffResult serialised as JSON for historical replay.
	DiffJSON  string    `json:"diff_json" gorm:"type:text"`
	CreatedAt time.Time `json:"created_at"`
}

// NewJoinerDetail is a staging table row for full new-joiner data uploaded separately.
type NewJoinerDetail struct {
	ID                int        `gorm:"primaryKey;autoIncrement" json:"id"`
	SubmissionID      int        `json:"submission_id" gorm:"not null;index"`
	MemberKey         string     `json:"member_key" gorm:"not null"`
	MemberName        string     `json:"member_name" gorm:"default:''"`
	EmployeeNumber    string     `json:"employee_number" gorm:"default:''"`
	IDNumber          string     `json:"id_number" gorm:"default:''"`
	MemberIdType      string     `json:"member_id_type" gorm:"default:''"`
	Gender            string     `json:"gender" gorm:"default:''"`
	DateOfBirth       string     `json:"date_of_birth" gorm:"default:''"`
	AnnualSalary      float64    `json:"annual_salary" gorm:"default:0"`
	SchemeCategory    string     `json:"scheme_category" gorm:"default:''"`
	AddressLine1      string     `json:"address_line_1" gorm:"default:''"`
	AddressLine2      string     `json:"address_line_2" gorm:"default:''"`
	City              string     `json:"city" gorm:"default:''"`
	Province          string     `json:"province" gorm:"default:''"`
	PostalCode        string     `json:"postal_code" gorm:"default:''"`
	PhoneNumber       string     `json:"phone_number" gorm:"default:''"`
	Email             string     `json:"email" gorm:"default:''"`
	Occupation        string     `json:"occupation" gorm:"default:''"`
	OccupationalClass string     `json:"occupational_class" gorm:"default:''"`
	EntryDate         string     `json:"entry_date" gorm:"default:''"`
	ValidationStatus  string     `json:"validation_status" gorm:"default:'valid'"`
	ExclusionReason   string     `json:"exclusion_reason" gorm:"default:''"`
	SyncedAt          *time.Time `json:"synced_at"`
	CreatedAt         time.Time  `json:"created_at"`
}
