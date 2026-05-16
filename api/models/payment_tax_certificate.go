package models

import "time"

// PaymentTaxCertificate is the audit + retrieval record for a tax-disclosure
// certificate (IT3(a) in SA) generated against a paid claim line. The actual
// rendered file lives on disk at StoragePath; this row tracks who/when and
// gives a stable download URL via certificate id.
//
// Phase 4 generates HTML certificates because there is no PDF library in
// the project yet. When a real renderer lands (gofpdf / unipdf / wkhtmltopdf),
// it produces the same content at the same StoragePath under a different
// extension — clients just request the certificate by id.
// Table name: payment_tax_certificates
type PaymentTaxCertificate struct {
	ID               int       `json:"id" gorm:"primaryKey;autoIncrement"`
	ScheduleID       int       `json:"schedule_id" gorm:"index"`
	ScheduleItemID   int       `json:"schedule_item_id" gorm:"index;uniqueIndex:idx_tax_cert_item_year"`
	ClaimID          int       `json:"claim_id" gorm:"index"`
	ClaimNumber      string    `json:"claim_number"`
	BenefitName      string    `json:"benefit_name"`
	BeneficiaryName  string    `json:"beneficiary_name"`
	BeneficiaryIDNumber string `json:"beneficiary_id_number"`
	TaxYear          int       `json:"tax_year" gorm:"uniqueIndex:idx_tax_cert_item_year"`
	GrossAmount      float64   `json:"gross_amount"`
	TaxWithheld      float64   `json:"tax_withheld"`
	CertificateRef   string    `json:"certificate_ref" gorm:"size:64;uniqueIndex"`
	FileName         string    `json:"file_name"`
	StoragePath      string    `json:"storage_path"`
	ContentType      string    `json:"content_type" gorm:"size:64"`
	GeneratedBy      string    `json:"generated_by"`
	GeneratedAt      time.Time `json:"generated_at" gorm:"autoCreateTime"`
}
