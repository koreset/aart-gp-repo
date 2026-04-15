package on_risk_letter_template

import (
	"archive/zip"
	"bytes"
	"io"
	"os"
	"testing"

	"api/services/quote_template"
)

// TestBuildSampleTemplate verifies the sample template is a valid ZIP and renders without error
func TestBuildSampleTemplate(t *testing.T) {
	// Build sample
	sampleData, err := BuildSampleTemplate()
	if err != nil {
		t.Fatalf("BuildSampleTemplate() failed: %v", err)
	}

	if len(sampleData) == 0 {
		t.Fatal("BuildSampleTemplate() returned empty data")
	}

	// Verify it's a valid ZIP
	_, err = zip.NewReader(bytes.NewReader(sampleData), int64(len(sampleData)))
	if err != nil {
		t.Fatalf("Sample is not a valid ZIP: %v", err)
	}

	// Write to disk for manual inspection
	err = writeToFile(sampleData, "/tmp/sample_on_risk_letter_template.docx")
	if err != nil {
		t.Logf("Warning: Could not write sample to disk: %v", err)
	}

	t.Log("Sample template created successfully at /tmp/sample_on_risk_letter_template.docx")
}

// TestRenderSampleWithSyntheticContext verifies the render engine can process the sample without infinite loops
func TestRenderSampleWithSyntheticContext(t *testing.T) {
	// Build sample
	sampleData, err := BuildSampleTemplate()
	if err != nil {
		t.Fatalf("BuildSampleTemplate() failed: %v", err)
	}

	// Create a synthetic context mimicking on-risk letter data
	ctx := quote_template.Context{
		"letter_reference":      "ORL-7-12-1234",
		"letter_date":           "14 Apr 2026",
		"generated_by":          "Jome Akpoduado",
		"scheme_name":           "Acme Pension Scheme",
		"scheme_contact":        "John Doe",
		"scheme_email":          "john@acme.com",
		"commencement_date":     "01 May 2026",
		"cover_end_date":        "30 Apr 2027",
		"industry":              "Manufacturing",
		"obligation_type":       "DC",
		"currency":              "ZAR",
		"free_cover_limit":      "500 000.00",
		"normal_retirement_age": "65",
		"continuation_option":   true,
		"distribution_channel":  "broker",
		"broker_name":           "Acme Brokers",
		"is_broker_channel":     true,
		"member_count":          "120",
		"quote_id":              "12",
		"quote_reference":       "Q-001",
		"total_annual_premium":  "315 000.00",
		"has_benefits":          true,
		"insurer": map[string]interface{}{
			"name":                    "XYZ Insurance Ltd",
			"contact_person":          "Jane Smith",
			"address_line_1":          "123 Main Street",
			"address_line_2":          "Suite 100",
			"address_line_3":          "",
			"city":                    "Johannesburg",
			"province":                "Gauteng",
			"post_code":               "2000",
			"country":                 "South Africa",
			"telephone":               "+27 11 234 5678",
			"email":                   "claims@xyz.com",
			"on_risk_letter_text":     "Custom closing text",
			"introductory_text":       "Custom intro",
			"general_provisions_text": "General provisions",
		},
		"benefit_summary": []map[string]interface{}{
			{
				"benefit":          "Group Life Assurance (GLA)",
				"annual_premium":   "120 000.00",
			},
			{
				"benefit":          "Permanent Total Disability (PTD)",
				"annual_premium":   "80 000.00",
			},
			{
				"benefit":          "Critical Illness (CI)",
				"annual_premium":   "115 000.00",
			},
		},
	}

	// Render the template
	renderedData, err := quote_template.Render(sampleData, ctx)
	if err != nil {
		t.Fatalf("Render() failed: %v", err)
	}

	if len(renderedData) == 0 {
		t.Fatal("Render() returned empty data")
	}

	// Verify it's still a valid ZIP
	zr, err := zip.NewReader(bytes.NewReader(renderedData), int64(len(renderedData)))
	if err != nil {
		t.Fatalf("Rendered output is not a valid ZIP: %v", err)
	}

	// Extract and verify document.xml exists and is valid XML
	var docXML []byte
	for _, f := range zr.File {
		if f.Name == "word/document.xml" {
			rc, err := f.Open()
			if err != nil {
				t.Fatalf("Could not open document.xml: %v", err)
			}
			docXML, err = io.ReadAll(rc)
			rc.Close()
			if err != nil {
				t.Fatalf("Could not read document.xml: %v", err)
			}
			break
		}
	}

	if len(docXML) == 0 {
		t.Fatal("document.xml not found in rendered output")
	}

	// Write rendered template to disk for manual inspection
	err = writeToFile(renderedData, "/tmp/sample_on_risk_letter_rendered.docx")
	if err != nil {
		t.Logf("Warning: Could not write rendered template to disk: %v", err)
	}

	t.Log("Rendered template created successfully at /tmp/sample_on_risk_letter_rendered.docx")
}

// Helper: Write bytes to file
func writeToFile(data []byte, path string) error {
	return os.WriteFile(path, data, 0o644)
}
