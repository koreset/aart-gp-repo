package services

import (
	"api/models"
	"testing"
	"time"
)

// calculateAgeLastBirthday must match the Excel formula
//   ROUNDDOWN((12*(YEAR(CommenDate)-YEAR(DoB)) + (MONTH(CommenDate)-MONTH(DoB)))/12, 0)
// which treats the birthday as "occurred" once the commencement month has
// reached the DoB month, with no day-of-month adjustment.
func TestCalculateAgeLastBirthday(t *testing.T) {
	utc := time.UTC
	mk := func(y int, m time.Month, d int) time.Time {
		return time.Date(y, m, d, 0, 0, 0, 0, utc)
	}

	cases := []struct {
		name        string
		commen, dob time.Time
		want        int
	}{
		{
			name:   "commencement before DoB month",
			commen: mk(2026, time.May, 1),
			dob:    mk(2000, time.June, 15),
			want:   25, // 12*26 + (5-6) = 311; 311/12 = 25
		},
		{
			name:   "commencement same month earlier day",
			commen: mk(2026, time.June, 1),
			dob:    mk(2000, time.June, 15),
			want:   26, // 12*26 + 0 = 312; 312/12 = 26 — no day-of-month adjustment
		},
		{
			name:   "commencement same month same day",
			commen: mk(2026, time.June, 15),
			dob:    mk(2000, time.June, 15),
			want:   26,
		},
		{
			name:   "commencement same month later day",
			commen: mk(2026, time.June, 30),
			dob:    mk(2000, time.June, 1),
			want:   26,
		},
		{
			name:   "commencement after DoB month",
			commen: mk(2026, time.July, 1),
			dob:    mk(2000, time.June, 15),
			want:   26, // 12*26 + 1 = 313; 313/12 = 26
		},
		{
			name:   "leap-year DoB Feb 29 - commencement Feb non-leap",
			commen: mk(2025, time.February, 28),
			dob:    mk(2000, time.February, 29),
			want:   25, // 12*25 + 0 = 300; 300/12 = 25
		},
		{
			name:   "exact year boundary January",
			commen: mk(2026, time.January, 1),
			dob:    mk(2000, time.January, 1),
			want:   26,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := calculateAgeLastBirthday(tc.commen, tc.dob)
			if got != tc.want {
				t.Fatalf("calculateAgeLastBirthday(%s, %s) = %d, want %d",
					tc.commen.Format("2006-01-02"), tc.dob.Format("2006-01-02"), got, tc.want)
			}
		})
	}
}

func TestCalculateMemberAge_Dispatcher(t *testing.T) {
	utc := time.UTC
	commen := time.Date(2026, time.May, 1, 0, 0, 0, 0, utc)
	dob := time.Date(2000, time.June, 15, 0, 0, 0, 0, utc)

	gotANB := calculateMemberAge(commen, dob, models.AgeMethodAgeNextBirthday)
	wantANB := calculateAgeNextBirthday(commen, dob)
	if gotANB != wantANB {
		t.Fatalf("calculateMemberAge ANB = %d, want %d", gotANB, wantANB)
	}

	gotALB := calculateMemberAge(commen, dob, models.AgeMethodAgeLastBirthday)
	wantALB := calculateAgeLastBirthday(commen, dob)
	if gotALB != wantALB {
		t.Fatalf("calculateMemberAge ALB = %d, want %d", gotALB, wantALB)
	}

	// Unknown method must fall back to ANB so a misconfigured row does not
	// produce a zero / negative age in downstream rate lookups.
	gotFallback := calculateMemberAge(commen, dob, "garbage_value")
	if gotFallback != wantANB {
		t.Fatalf("calculateMemberAge fallback for unknown method = %d, want %d (ANB)",
			gotFallback, wantANB)
	}
}
