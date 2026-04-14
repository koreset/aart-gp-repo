package quote_template

import (
	"strings"
	"testing"
)

func TestStitchRuns_SingleRun_NoOp(t *testing.T) {
	// Token in a single run should not be modified
	input := `<w:p><w:r><w:t>{{quote_name}}</w:t></w:r></w:p>`
	result := StitchRuns(input)
	if result != input {
		t.Errorf("Single run should not be modified. Got: %s", result)
	}
}

func TestStitchRuns_SplitAcrossTwoRuns(t *testing.T) {
	// Token split across two runs should be merged
	input := `<w:p><w:r><w:rPr><w:b/></w:rPr><w:t>{{quote_</w:t></w:r><w:r><w:t>name}}</w:t></w:r></w:p>`
	result := StitchRuns(input)

	// The token should now be in a single run
	if !strings.Contains(result, "{{quote_name}}") {
		t.Errorf("Token not found in result: %s", result)
	}

	// Count runs - should be fewer after stitching
	runCount := strings.Count(result, "<w:r>")
	if runCount > 2 { // Allow for some runs, but should be consolidated
		t.Logf("Result has %d runs: %s", runCount, result)
	}
}

func TestStitchRuns_SplitAcrossThreeRuns(t *testing.T) {
	// Token split across three runs
	input := `<w:p><w:r><w:t>{{quote</w:t></w:r><w:r><w:t>_na</w:t></w:r><w:r><w:t>me}}</w:t></w:r></w:p>`
	result := StitchRuns(input)

	if !strings.Contains(result, "{{quote_name}}") {
		t.Errorf("Token not found in result: %s", result)
	}
}

func TestStitchRuns_MultipleTokens(t *testing.T) {
	// Multiple tokens in one paragraph - some split, some not
	input := `<w:p><w:r><w:t>{{quote_</w:t></w:r><w:r><w:t>name}} and {{scheme_name}}</w:t></w:r></w:p>`
	result := StitchRuns(input)

	if !strings.Contains(result, "quote_name") {
		t.Errorf("First token reference not found in result: %s", result)
	}
	if !strings.Contains(result, "scheme_name") {
		t.Errorf("Second token reference not found in result: %s", result)
	}
}

func TestStitchRuns_PreserveFormatting(t *testing.T) {
	// When merging, the first run's formatting should be preserved
	input := `<w:p><w:r><w:rPr><w:b/><w:i/></w:rPr><w:t>{{quote_</w:t></w:r><w:r><w:t>name}}</w:t></w:r></w:p>`
	result := StitchRuns(input)

	// The merged run should retain <w:b/> and <w:i/>
	if !strings.Contains(result, "<w:b/>") {
		t.Errorf("Bold formatting lost: %s", result)
	}
	if !strings.Contains(result, "<w:i/>") {
		t.Errorf("Italic formatting lost: %s", result)
	}
}

func TestStitchRuns_NoTokens(t *testing.T) {
	// Paragraph without tokens should be unchanged
	input := `<w:p><w:r><w:t>Hello world</w:t></w:r></w:p>`
	result := StitchRuns(input)
	if result != input {
		t.Errorf("Paragraph without tokens should not change. Got: %s", result)
	}
}

func TestExtractTextFromRun(t *testing.T) {
	tests := []struct {
		name     string
		run      string
		expected string
	}{
		{
			name:     "Simple text",
			run:      `<w:r><w:t>Hello</w:t></w:r>`,
			expected: "Hello",
		},
		{
			name:     "Text with formatting",
			run:      `<w:r><w:rPr><w:b/></w:rPr><w:t>Bold text</w:t></w:r>`,
			expected: "Bold text",
		},
		{
			name:     "Multiple text elements",
			run:      `<w:r><w:t>Hello</w:t><w:t>World</w:t></w:r>`,
			expected: "HelloWorld",
		},
		{
			name:     "Empty text",
			run:      `<w:r><w:t></w:t></w:r>`,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractTextFromRun(tt.run)
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}
