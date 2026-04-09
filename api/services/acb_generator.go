package services

import (
	"api/models"
	"fmt"
	"math"
	"regexp"
	"strings"
	"time"
)

// Universal branch codes for SA banks (used in ACB files).
var universalBranchCodes = map[string]string{
	"FNB":            "250655",
	"Standard Bank":  "051001",
	"ABSA":           "632005",
	"Nedbank":        "198765",
	"Capitec":        "470010",
	"Investec":       "580105",
	"African Bank":   "430000",
	"TymeBank":       "678910",
	"Discovery Bank": "679000",
	"Bank Zero":      "888000",
}

// GetUniversalBranchCode returns the universal branch code for a given bank name.
func GetUniversalBranchCode(bankName string) string {
	if code, ok := universalBranchCodes[bankName]; ok {
		return code
	}
	return ""
}

// padRight left-justifies s and pads with spaces to the given length.
func padRight(s string, length int) string {
	if len(s) >= length {
		return s[:length]
	}
	return s + strings.Repeat(" ", length-len(s))
}

// padLeft right-justifies s and pads with zeros to the given length.
func padLeft(s string, length int) string {
	if len(s) >= length {
		return s[len(s)-length:]
	}
	return strings.Repeat("0", length-len(s)) + s
}

// amountToCents converts a rand amount to cents.
func amountToCents(amount float64) int64 {
	return int64(math.Round(amount * 100))
}

// formatCents formats cents as a zero-padded string of given width.
func formatCents(cents int64, width int) string {
	return padLeft(fmt.Sprintf("%d", cents), width)
}

// formatDateYYMMDD formats a time.Time as YYMMDD.
func formatDateYYMMDD(t time.Time) string {
	return t.Format("060102")
}

// formatDateCCYYMMDD formats a time.Time as CCYYMMDD.
func formatDateCCYYMMDD(t time.Time) string {
	return t.Format("20060102")
}

// sanitizeASCII strips non-ASCII printable characters from a string.
func sanitizeASCII(s string) string {
	re := regexp.MustCompile(`[^\x20-\x7E]`)
	return re.ReplaceAllString(s, "")
}

// buildRecord concatenates fields and pads/truncates to exactly 179 chars, then appends CR+LF.
func buildRecord(fields ...string) string {
	line := strings.Join(fields, "")
	if len(line) > 179 {
		line = line[:179]
	} else if len(line) < 179 {
		line = line + strings.Repeat(" ", 179-len(line))
	}
	return line + "\r\n"
}

// accountToNumeric converts an account number string to int64 for hash total calculation.
func accountToNumeric(acct string) int64 {
	var result int64
	for _, ch := range acct {
		if ch >= '0' && ch <= '9' {
			result = result*10 + int64(ch-'0')
		}
	}
	return result
}

// transactionCode returns the ACB transaction code based on service type.
func transactionCode(serviceType string) string {
	switch serviceType {
	case "same_day":
		return "040000"
	case "one_day":
		return "010000"
	default: // two_day
		return "020000"
	}
}

// purgeDateFromCreation returns creation date + 5 business days.
func purgeDateFromCreation(creation time.Time) time.Time {
	days := 0
	t := creation
	for days < 5 {
		t = t.AddDate(0, 0, 1)
		if t.Weekday() != time.Saturday && t.Weekday() != time.Sunday {
			days++
		}
	}
	return t
}

// GenerateACBFileContent builds a complete ACB file from profile and payment items.
// Returns the file bytes, hash total, and any error.
func GenerateACBFileContent(profile models.ACBBankProfile, items []models.ClaimPaymentScheduleItem, actionDate time.Time) ([]byte, int64, error) {
	if len(items) == 0 {
		return nil, 0, fmt.Errorf("no items to generate ACB file for")
	}

	creationDate := time.Now()
	purgeDate := purgeDateFromCreation(creationDate)
	var records []string

	// Record Type 02 — Installation Header
	rec02 := buildRecord(
		"02",                                        // pos 1-2
		"    ",                                      // pos 3-6 filler
		padRight("MAGTAPE", 10),                     // pos 7-16
		padRight(formatDateCCYYMMDD(creationDate), 10), // pos 17-26
		formatDateCCYYMMDD(purgeDate)[:6],           // pos 27-32 (CCYYMM portion — 6 chars)
	)
	records = append(records, rec02)

	// Record Type 04 — User Header
	rec04 := buildRecord(
		"04",                                        // pos 1-2
		padRight(profile.UserCode, 4),               // pos 3-6
		padRight(formatDateCCYYMMDD(creationDate), 10), // pos 7-16
		formatDateCCYYMMDD(purgeDate)[:6],           // pos 17-22
		formatDateYYMMDD(actionDate),                // pos 23-28 first action date
		formatDateYYMMDD(actionDate),                // pos 29-34 last action date
		"0001",                                      // pos 35-38 first sequence number
		padLeft(fmt.Sprintf("%d", profile.GenerationNumber), 4), // pos 39-42
		padRight(profile.BankTypeCode, 2),           // pos 43-44
	)
	records = append(records, rec04)

	// Record Type 50 — Credit Transactions
	var hashTotal int64
	var totalCreditCents int64
	creditCount := 0

	for _, item := range items {
		cents := amountToCents(item.ClaimAmount)
		totalCreditCents += cents

		destBranch := item.BankBranchCode
		if destBranch == "" {
			destBranch = GetUniversalBranchCode(item.BankName)
		}

		acctNum := padLeft(item.BankAccountNumber, 11)
		hashTotal += accountToNumeric(acctNum)

		acctType := item.BankAccountType
		if acctType == "" {
			acctType = "1"
		}

		holderName := sanitizeASCII(item.AccountHolderName)
		if holderName == "" {
			holderName = sanitizeASCII(item.MemberName)
		}

		// Build user reference: schedule item ID + claim number
		userRef := fmt.Sprintf("%d-%s", item.ScheduleID, item.ClaimNumber)

		rec50 := buildRecord(
			"50",                           // pos 1-2
			padLeft(destBranch, 6),          // pos 3-8
			acctNum,                         // pos 9-19
			padRight(acctType, 1),           // pos 20
			transactionCode(profile.ServiceType), // pos 21-26
			formatCents(cents, 12),          // pos 27-38
			formatDateYYMMDD(actionDate),    // pos 39-44
			"3",                             // pos 45 entry class (ACB credit)
			"00",                            // pos 46-47 tax code
			padRight(userRef, 30),           // pos 48-77
			"0000000000",                    // pos 78-87 homing institution
			padRight(holderName, 20),        // pos 88-107
		)
		records = append(records, rec50)
		creditCount++
	}

	// Hash total mod 10^12
	hashTotal = hashTotal % 1000000000000

	// Record Type 52 — Contra (debit from source account)
	sourceAcct := padLeft(profile.UserAccountNumber, 11)
	sourceType := profile.UserAccountType
	if sourceType == "" {
		sourceType = "1"
	}

	rec52 := buildRecord(
		"52",                                  // pos 1-2
		padLeft(profile.UserBranchCode, 6),    // pos 3-8
		sourceAcct,                            // pos 9-19
		padRight(sourceType, 1),               // pos 20
		transactionCode(profile.ServiceType),  // pos 21-26
		formatCents(totalCreditCents, 12),     // pos 27-38
		formatDateYYMMDD(actionDate),          // pos 39-44
		"3",                                   // pos 45
		"00",                                  // pos 46-47
		padRight("ACB CONTRA", 30),            // pos 48-77
		"0000000000",                          // pos 78-87
		padRight(profile.ProfileName, 20),     // pos 88-107
	)
	records = append(records, rec52)
	debitCount := 1

	// Record Type 92 — User Trailer
	totalDataRecords := creditCount + debitCount
	rec92 := buildRecord(
		"92",                                                          // pos 1-2
		padRight(profile.UserCode, 4),                                 // pos 3-6
		padLeft("1", 12),                                              // pos 7-18 first seq
		padLeft(fmt.Sprintf("%d", totalDataRecords), 12),              // pos 19-30 last seq
		padLeft(fmt.Sprintf("%d", hashTotal), 12),                     // pos 31-42 hash total
		padLeft(fmt.Sprintf("%d", debitCount), 12),                    // pos 43-54 debit records
		padLeft(fmt.Sprintf("%d", creditCount), 12),                   // pos 55-66 credit records
		formatCents(totalCreditCents, 12),                             // pos 67-78 total debit (cents)
		formatCents(totalCreditCents, 12),                             // pos 79-90 total credit (cents)
	)
	records = append(records, rec92)

	// Record Type 94 — Installation Trailer
	totalRecordCount := len(records) + 1 // include this record itself
	rec94 := buildRecord(
		"94",                                                   // pos 1-2
		"    ",                                                 // pos 3-6 filler
		padLeft(fmt.Sprintf("%d", totalRecordCount), 12),       // pos 7-18
	)
	records = append(records, rec94)

	content := strings.Join(records, "")
	return []byte(content), hashTotal, nil
}
