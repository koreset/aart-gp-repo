package quote_template

import (
	"regexp"
	"strings"
)

// StitchRuns processes Word document XML to merge fragmented runs that contain template tokens.
// When Word saves a document with curly braces, it may split tokens across multiple <w:r> elements.
// This function finds token patterns ({{...}}) that span multiple runs and merges them into a single run.
func StitchRuns(documentXML string) string {
	// Find all paragraphs
	paraRegex := regexp.MustCompile(`<w:p>.*?</w:p>`)
	return paraRegex.ReplaceAllStringFunc(documentXML, stitchParagraphRuns)
}

// runInfo represents a single run with its position in the flat text
type runInfo struct {
	index int
	text  string
	run   string // full run element
	start int    // character position in concatenated text
	end   int    // character position in concatenated text
}

// stitchParagraphRuns processes a single paragraph to merge runs containing split tokens
func stitchParagraphRuns(para string) string {
	// Extract runs from the paragraph
	runRegex := regexp.MustCompile(`<w:r>.*?</w:r>`)
	runs := runRegex.FindAllString(para, -1)

	if len(runs) == 0 {
		return para
	}

	var runList []runInfo
	totalText := strings.Builder{}

	for i, run := range runs {
		text := extractTextFromRun(run)
		start := totalText.Len()
		totalText.WriteString(text)
		end := totalText.Len()

		runList = append(runList, runInfo{
			index: i,
			text:  text,
			run:   run,
			start: start,
			end:   end,
		})
	}

	flatText := totalText.String()

	// Find all tokens in the flat text
	tokenRegex := regexp.MustCompile(`\{\{[^}]*\}\}`)
	tokens := tokenRegex.FindAllStringIndex(flatText, -1)

	// For each token that spans multiple runs, mark those runs for merging
	toMerge := make(map[int]bool) // maps from run index
	mergeGroups := make(map[int]int)  // maps from run index to group ID

	for _, tokenRange := range tokens {
		tokenStart := tokenRange[0]
		tokenEnd := tokenRange[1]

		// Find which runs this token spans
		var startRunIdx, endRunIdx int
		for i, ri := range runList {
			if ri.start <= tokenStart && tokenStart < ri.end {
				startRunIdx = i
			}
			if ri.start < tokenEnd && tokenEnd <= ri.end {
				endRunIdx = i
			}
		}

		// If token spans multiple runs, merge them
		if startRunIdx != endRunIdx {
			for i := startRunIdx; i <= endRunIdx; i++ {
				toMerge[i] = true
				mergeGroups[i] = startRunIdx
			}
		}
	}

	// If no merges needed, return original paragraph
	if len(toMerge) == 0 {
		return para
	}

	// Rebuild the paragraph by merging identified runs
	result := strings.Builder{}
	i := 0

	for i < len(runList) {
		if !toMerge[i] {
			// This run is not part of a merge group, keep it as-is
			result.WriteString(runList[i].run)
			i++
		} else {
			// This run starts a merge group
			groupID := mergeGroups[i]
			if i != groupID {
				// We're not the start of the group, skip
				i++
				continue
			}

			// Collect all runs in this merge group
			var mergeRuns []runInfo
			j := i
			for j < len(runList) && toMerge[j] && mergeGroups[j] == groupID {
				mergeRuns = append(mergeRuns, runList[j])
				j++
			}

			// Merge these runs: take the first run's formatting, update its text
			if len(mergeRuns) > 0 {
				mergedRun := mergeRunsIntoOne(mergeRuns)
				result.WriteString(mergedRun)
			}

			i = j
		}
	}

	// Rebuild the full paragraph
	return rebuildParagraph(para, result.String())
}

// extractTextFromRun extracts the text content from a run element
func extractTextFromRun(run string) string {
	// Find <w:t>...</w:t> tags
	textRegex := regexp.MustCompile(`<w:t[^>]*>([^<]*)</w:t>`)
	matches := textRegex.FindAllStringSubmatch(run, -1)

	var text strings.Builder
	for _, match := range matches {
		if len(match) > 1 {
			text.WriteString(match[1])
		}
	}

	return text.String()
}

// mergeRunsIntoOne merges multiple runs into a single run, preserving the first run's formatting
func mergeRunsIntoOne(runs []runInfo) string {
	if len(runs) == 0 {
		return ""
	}

	firstRun := runs[0].run
	var totalText strings.Builder
	for _, ri := range runs {
		totalText.WriteString(ri.text)
	}

	// Extract the run properties from the first run
	var rPr string
	rPrRegex := regexp.MustCompile(`<w:rPr>.*?</w:rPr>`)
	rPrMatches := rPrRegex.FindString(firstRun)
	if rPrMatches != "" {
		rPr = rPrMatches
	}

	// Extract any other attributes from first run's <w:r> tag
	var openTag string
	openRegex := regexp.MustCompile(`<w:r[^>]*>`)
	openMatches := openRegex.FindString(firstRun)
	if openMatches != "" {
		openTag = openMatches
	} else {
		openTag = `<w:r>`
	}

	// Build the merged run
	mergedText := totalText.String()
	var result strings.Builder
	result.WriteString(openTag)
	if rPr != "" {
		result.WriteString(rPr)
	}
	result.WriteString(`<w:t>`)
	result.WriteString(mergedText)
	result.WriteString(`</w:t>`)
	result.WriteString(`</w:r>`)

	return result.String()
}

// rebuildParagraph replaces the runs section of a paragraph with new runs
func rebuildParagraph(originalPara, newRuns string) string {
	// Extract paragraph properties (if any)
	pPrRegex := regexp.MustCompile(`<w:p(?: [^>]*)?>(?:<w:pPr>.*?</w:pPr>)?`)
	pPrMatch := pPrRegex.FindString(originalPara)

	// Extract paragraph end tag
	endRegex := regexp.MustCompile(`</w:p>`)
	endMatch := endRegex.FindString(originalPara)

	if pPrMatch == "" {
		pPrMatch = `<w:p>`
	}

	if endMatch == "" {
		endMatch = `</w:p>`
	}

	return pPrMatch + newRuns + endMatch
}
