package quote_template

import (
	"archive/zip"
	"bytes"
	"os"
	"testing"
)

func TestE2E_SampleAndRender(t *testing.T) {
	sample, err := BuildSampleTemplate()
	if err != nil {
		t.Fatalf("BuildSampleTemplate: %v", err)
	}
	if _, err := zip.NewReader(bytes.NewReader(sample), int64(len(sample))); err != nil {
		t.Fatalf("sample is not a valid ZIP: %v", err)
	}
	os.WriteFile("/tmp/sample_quote_template.docx", sample, 0o644)
	t.Logf("sample template: %d bytes", len(sample))

	ctx := Context{
		"quote_name":               "Q-DEMO-001",
		"quote_number":             "Q-DEMO-001",
		"scheme_name":              "Demo Scheme",
		"creation_date":            "14 Apr 2026",
		"commencement_date":        "01 May 2026",
		"industry":                 "Manufacturing",
		"currency":                 "ZAR",
		"free_cover_limit":         "500 000.00",
		"total_lives":              "120",
		"total_sum_assured":        "84 000 000.00",
		"total_annual_salary":      "28 000 000.00",
		"total_annual_premium":     "315 000.00",
		"has_non_funeral_benefits": true,
		"insurer": map[string]any{
			"name":                    "Example Insurer Ltd",
			"address_line_1":          "1 Example Street",
			"address_line_2":          "Suite 200",
			"address_line_3":          "",
			"city":                    "Johannesburg",
			"province":                "Gauteng",
			"post_code":               "2000",
			"telephone":               "+27 11 000 0000",
			"email":                   "quotes@example.com",
			"introductory_text":       "We're pleased to submit this quote for your consideration.",
			"general_provisions_text": "Standard underwriting provisions and exclusions apply.",
		},
		"categories": []any{
			map[string]any{
				"name":              "Management",
				"member_count":      "20",
				"total_salary":      "10 000 000.00",
				"total_sum_assured": "30 000 000.00",
				"annual_premium":    "120 000.00",
				"percent_salary":    "1.20%",
				"has_gla":           true,
				"has_funeral":       true,
				"gla": map[string]any{
					"title":             "Group Life",
					"total_sum_assured": "30 000 000.00",
					"annual_premium":    "120 000.00",
				},
				"funeral": map[string]any{
					"total_annual_premium": "9 000.00",
				},
			},
			map[string]any{
				"name":              "Staff",
				"member_count":      "100",
				"total_salary":      "18 000 000.00",
				"total_sum_assured": "54 000 000.00",
				"annual_premium":    "195 000.00",
				"percent_salary":    "1.08%",
				"has_gla":           true,
				"has_funeral":       true,
				"gla": map[string]any{
					"title":             "Group Life",
					"total_sum_assured": "54 000 000.00",
					"annual_premium":    "195 000.00",
				},
				"funeral": map[string]any{
					"total_annual_premium": "18 000.00",
				},
			},
		},
	}

	rendered, err := Render(sample, ctx)
	if err != nil {
		t.Fatalf("Render: %v", err)
	}
	if _, err := zip.NewReader(bytes.NewReader(rendered), int64(len(rendered))); err != nil {
		t.Fatalf("rendered is not a valid ZIP: %v", err)
	}
	os.WriteFile("/tmp/sample_quote_rendered.docx", rendered, 0o644)
	t.Logf("rendered: %d bytes", len(rendered))
}
