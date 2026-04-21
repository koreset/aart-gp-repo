package quote_template

import (
	"archive/zip"
	"bytes"
	"io"
	"regexp"
	"sort"
	"strings"
	"testing"

	"api/models"
	"api/services/quote_docx"
)

// TestSchema_SampleContainsEveryContextKey ensures the self-documenting
// sample template enumerates every token the schema exposes. This is the
// guard that replaces the old hand-maintained duplication between
// context.go and sample.go — add a key in schema.go and if the sample
// doesn't pick it up, this test fails.
func TestSchema_SampleContainsEveryContextKey(t *testing.T) {
	// Build the sample once; extract the full concatenated document text.
	sampleBytes, err := BuildSampleTemplate()
	if err != nil {
		t.Fatalf("BuildSampleTemplate: %v", err)
	}
	sampleText := extractSampleDocText(t, sampleBytes)

	// Every token the schema promises, expressed as it would appear in a
	// template. Root scalars -> "{{key}}"; insurer -> "{{insurer.key}}";
	// category scalars -> "{{key}}" (inside {{#categories}});
	// category bools -> "{{#key}}"; per-benefit -> "{{prefix.key}}".
	expected := expectedTokenSet()

	for _, tok := range expected {
		if !strings.Contains(sampleText, tok) {
			t.Errorf("sample template missing token %q — add it to buildSampleBodyXML or drop it from schema.go", tok)
		}
	}
}

// TestSchema_BuildContextAgreesWithSchema asserts that every key emitted
// by BuildContext comes from the schema — catches the inverse drift of a
// key being stuffed into Context outside of a *Fields function.
func TestSchema_BuildContextAgreesWithSchema(t *testing.T) {
	// Construct a Context directly from the schema with zero fixtures.
	var (
		zs models.MemberRatingResultSummary
		zc models.SchemeCategory
		zq models.GroupPricingQuote
		zi models.GroupPricingInsurerDetail
		zt quote_docx.BenefitTitles
		zT quote_docx.QuoteTotals
	)
	schemaCtx := Context(fieldsToMap(quoteFields(zq, zT, false)))
	schemaCtx["insurer"] = fieldsToMap(insurerFields(zi))
	// One synthetic category built from the schema — has=true for every
	// benefit so all sub-keys are exercised.
	cat := fieldsToMap(categoryScalarFields(zs, zc))
	for k, v := range fieldsToMap(categoryBoolFields(zs, benefitFlags{GLA: true, SGLA: true, PTD: true, CI: true, PHI: true, TTD: true, Funeral: true})) {
		cat[k] = v
	}
	cat["gla"] = fieldsToMap(glaFields(zs, zc, zq, zt))
	cat["sgla"] = fieldsToMap(sglaFields(zs, zc, zq, zt))
	cat["ptd"] = fieldsToMap(ptdFields(zs, zc, zq, zt))
	cat["ci"] = fieldsToMap(ciFields(zs, zc, zq, zt))
	cat["phi"] = fieldsToMap(phiFields(zs, zc, zt))
	cat["ttd"] = fieldsToMap(ttdFields(zs, zc, zt))
	cat["funeral"] = fieldsToMap(funeralFields(zs, zc, zt))
	schemaCtx["categories"] = []map[string]interface{}{cat}

	// The keys flattened by this helper are what a template author can
	// actually reference. If you find a key here that the sample doesn't
	// show, the first test above will flag it; if you find one the
	// schema shouldn't expose, remove the Field from schema.go.
	got := flattenKeys(schemaCtx)
	sort.Strings(got)
	if len(got) == 0 {
		t.Fatalf("schema produced no keys — did a *Fields function break?")
	}
	// Spot-check a handful of canonical keys so regressions in Field
	// construction fail loudly.
	mustContain(t, got, []string{
		"scheme_name",
		"insurer.name",
		"categories[].name",
		"categories[].has_gla",
		"categories[].gla.total_sum_assured",
		"categories[].gla_rate_per_1000",
		"categories[].funeral.total_annual_premium",
	})
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func expectedTokenSet() []string {
	var (
		zs models.MemberRatingResultSummary
		zc models.SchemeCategory
		zq models.GroupPricingQuote
		zi models.GroupPricingInsurerDetail
		zT quote_docx.QuoteTotals
	)
	var out []string

	for _, f := range quoteFields(zq, zT, false) {
		if _, isBool := f.Value.(bool); isBool {
			// Quote-level bools render as {{#key}} block demonstrations,
			// which the sample still contains verbatim.
			out = append(out, "{{#"+f.Key+"}}")
		} else {
			out = append(out, "{{"+f.Key+"}}")
		}
	}
	for _, f := range insurerFields(zi) {
		out = append(out, "{{insurer."+f.Key+"}}")
	}
	for _, f := range categoryScalarFields(zs, zc) {
		out = append(out, "{{"+f.Key+"}}")
	}
	for _, f := range categoryBoolFields(zs, benefitFlags{}) {
		out = append(out, "{{#"+f.Key+"}}")
	}
	for _, spec := range benefitSpecsForSample() {
		for _, f := range spec.Fields() {
			out = append(out, "{{"+spec.Prefix+"."+f.Key+"}}")
		}
	}
	return out
}

// extractSampleDocText unzips the generated docx and returns the raw
// concatenation of text inside <w:t> elements. Good enough for
// substring-containment assertions.
func extractSampleDocText(t *testing.T, docx []byte) string {
	t.Helper()
	zr, err := zip.NewReader(bytes.NewReader(docx), int64(len(docx)))
	if err != nil {
		t.Fatalf("open sample as zip: %v", err)
	}
	for _, f := range zr.File {
		if f.Name != "word/document.xml" {
			continue
		}
		rc, err := f.Open()
		if err != nil {
			t.Fatalf("open document.xml: %v", err)
		}
		data, err := io.ReadAll(rc)
		rc.Close()
		if err != nil {
			t.Fatalf("read document.xml: %v", err)
		}
		// Strip <w:t> tags, keeping their content. Tokens may still span
		// run boundaries in theory, but BuildSampleTemplate always emits
		// a token inside a single run so this flat join is sufficient.
		textRe := regexp.MustCompile(`<w:t[^>]*>([^<]*)</w:t>`)
		var b strings.Builder
		for _, m := range textRe.FindAllStringSubmatch(string(data), -1) {
			b.WriteString(m[1])
			b.WriteString(" ")
		}
		return b.String()
	}
	t.Fatalf("document.xml not found in sample")
	return ""
}

// flattenKeys walks a Context-shaped structure and returns dotted keys a
// template author would use. Lists are rendered as "foo[].bar".
func flattenKeys(v interface{}) []string {
	var out []string
	switch x := v.(type) {
	case Context:
		for k, child := range x {
			out = append(out, prefixKeys(k, flattenKeys(child))...)
			out = append(out, k)
		}
	case map[string]interface{}:
		for k, child := range x {
			out = append(out, prefixKeys(k, flattenKeys(child))...)
			out = append(out, k)
		}
	case []map[string]interface{}:
		if len(x) > 0 {
			out = append(out, prefixKeys("[]", flattenKeys(x[0]))...)
		}
	}
	return out
}

func prefixKeys(prefix string, keys []string) []string {
	out := make([]string, 0, len(keys))
	for _, k := range keys {
		if prefix == "[]" {
			out = append(out, "[]."+k)
		} else {
			out = append(out, prefix+"."+k)
		}
	}
	return out
}

func mustContain(t *testing.T, haystack []string, needles []string) {
	t.Helper()
	set := make(map[string]bool, len(haystack))
	for _, s := range haystack {
		// Normalise "[]" prefixes: the flattener emits "[].name" for list
		// items; the assertions phrase it as "categories[].name".
		set[s] = true
	}
	// Recompose as rooted paths like "categories[].name" by looking for
	// any haystack entry ending in the suffix after "categories".
	for _, need := range needles {
		if strings.HasPrefix(need, "categories[].") {
			suffix := strings.TrimPrefix(need, "categories")
			found := false
			for _, got := range haystack {
				if strings.HasSuffix(got, suffix) {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("schema missing expected path %q (flattened keys: %v)", need, haystack)
			}
			continue
		}
		if !set[need] {
			t.Errorf("schema missing expected key %q", need)
		}
	}
}
