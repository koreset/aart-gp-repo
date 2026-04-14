package quote_template

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strings"
)

// Render processes a template DOCX and substitutes tokens with context values
func Render(templateBytes []byte, ctx Context) ([]byte, error) {
	// Open template as ZIP
	reader := bytes.NewReader(templateBytes)
	zr, err := zip.NewReader(reader, int64(len(templateBytes)))
	if err != nil {
		return nil, fmt.Errorf("failed to open template as ZIP: %w", err)
	}

	// Collect all files from the ZIP
	var files map[string][]byte = make(map[string][]byte)
	var filesToWrite []string

	for _, f := range zr.File {
		rc, err := f.Open()
		if err != nil {
			return nil, fmt.Errorf("failed to read file %s: %w", f.Name, err)
		}
		data, err := io.ReadAll(rc)
		rc.Close()
		if err != nil {
			return nil, fmt.Errorf("failed to read content of %s: %w", f.Name, err)
		}
		files[f.Name] = data
		filesToWrite = append(filesToWrite, f.Name)
	}

	// Process document.xml and related XMLs
	xmlFiles := []string{
		"word/document.xml",
		"word/header1.xml",
		"word/header2.xml",
		"word/header3.xml",
		"word/footer1.xml",
		"word/footer2.xml",
		"word/footer3.xml",
	}

	for _, xmlFile := range xmlFiles {
		if data, exists := files[xmlFile]; exists {
			// Convert bytes to string
			xmlStr := string(data)

			// Stitch runs to merge fragmented tokens
			xmlStr = StitchRuns(xmlStr)

			// Perform substitution
			xmlStr = substituteTokens(xmlStr, ctx)

			// Convert back to bytes
			files[xmlFile] = []byte(xmlStr)
		}
	}

	// Build output ZIP
	outBuf := new(bytes.Buffer)
	zw := zip.NewWriter(outBuf)

	for _, filename := range filesToWrite {
		data := files[filename]
		w, err := zw.Create(filename)
		if err != nil {
			return nil, fmt.Errorf("failed to create ZIP entry %s: %w", filename, err)
		}
		_, err = w.Write(data)
		if err != nil {
			return nil, fmt.Errorf("failed to write ZIP entry %s: %w", filename, err)
		}
	}

	err = zw.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to close output ZIP: %w", err)
	}

	return outBuf.Bytes(), nil
}

// substituteTokens performs all token substitutions in an XML string
func substituteTokens(xmlStr string, ctx Context) string {
	// First pass: process conditional and iteration blocks
	xmlStr = processBlocks(xmlStr, ctx)

	// Second pass: simple substitutions
	xmlStr = performSimpleSubstitutions(xmlStr, ctx)

	return xmlStr
}

// processBlocks handles {{#key}}...{{/key}} blocks (both conditionals and iterations)
func processBlocks(xmlStr string, ctx Context) string {
	// Process blocks iteratively, handling innermost first
	for {
		block := findInnermostBlock(xmlStr)
		if block == nil {
			break
		}

		// Resolve the key in context
		value := resolvePath(ctx, block.key)

		// Process based on value type
		var replacement string
		switch v := value.(type) {
		case bool:
			if v {
				// Include body, remove conditional tags
				replacement = removeBlockMarkersSimple(block.body, block.key)
			} else {
				// Exclude body and the block markers
				replacement = ""
			}
		case []map[string]interface{}:
			// Iteration over array of objects
			var expanded strings.Builder
			for _, item := range v {
				// Create new context with item fields
				itemCtx := mergeContexts(ctx, item)
				// Process body with new context
				processedBody := processBlocks(block.body, itemCtx)
				processedBody = performSimpleSubstitutions(processedBody, itemCtx)
				expanded.WriteString(processedBody)
			}
			replacement = expanded.String()
		case []interface{}:
			// Iteration over generic array
			var expanded strings.Builder
			for _, item := range v {
				itemCtx := mergeContexts(ctx, map[string]interface{}{"_value": item})
				processedBody := processBlocks(block.body, itemCtx)
				processedBody = performSimpleSubstitutions(processedBody, itemCtx)
				expanded.WriteString(processedBody)
			}
			replacement = expanded.String()
		default:
			// Not a recognized block type, leave as-is
			replacement = xmlStr[block.start:block.end]
		}

		xmlStr = xmlStr[:block.start] + replacement + xmlStr[block.end:]
	}

	return xmlStr
}

// blockInfo represents a found block
type blockInfo struct {
	key   string
	body  string
	start int
	end   int
}

// findInnermostBlock finds the innermost complete block
func findInnermostBlock(xmlStr string) *blockInfo {
	// Simple greedy approach: find first {{ # and match with closest {{/
	openIdx := strings.Index(xmlStr, "{{#")
	if openIdx == -1 {
		return nil
	}

	// Extract the key
	keyStart := openIdx + 3
	keyEnd := strings.Index(xmlStr[keyStart:], "}}")
	if keyEnd == -1 {
		return nil
	}
	keyEnd += keyStart

	key := xmlStr[keyStart:keyEnd]
	bodyStart := keyEnd + 2

	// Find the matching close tag
	closeTag := "{{/" + key + "}}"
	closeIdx := strings.Index(xmlStr[bodyStart:], closeTag)
	if closeIdx == -1 {
		return nil
	}
	closeIdx += bodyStart

	return &blockInfo{
		key:   key,
		body:  xmlStr[bodyStart:closeIdx],
		start: openIdx,
		end:   closeIdx + len(closeTag),
	}
}

// removeBlockMarkersSimple removes {{#key}} and {{/key}} markers from text
func removeBlockMarkersSimple(body string, key string) string {
	openPattern := regexp.MustCompile(`\{\{#` + regexp.QuoteMeta(key) + `\}\}`)
	closePattern := regexp.MustCompile(`\{\{/` + regexp.QuoteMeta(key) + `\}\}`)

	body = openPattern.ReplaceAllString(body, "")
	body = closePattern.ReplaceAllString(body, "")

	return body
}


// removeParagraphIfOnlyContent removes an entire paragraph if the block marker is its only content
func removeParagraphIfOnlyContent(xmlStr string, blockStart, blockEnd int) string {
	// Find the containing paragraph
	paraOpenIdx := strings.LastIndex(xmlStr[:blockStart], "<w:p")
	if paraOpenIdx == -1 {
		return xmlStr
	}
	paraOpenClose := strings.Index(xmlStr[paraOpenIdx:], ">")
	if paraOpenClose == -1 {
		return xmlStr
	}
	paraOpenClose += paraOpenIdx

	paraCloseIdx := strings.Index(xmlStr[blockEnd:], "</w:p>")
	if paraCloseIdx == -1 {
		return xmlStr
	}
	paraCloseIdx += blockEnd

	// Check if the paragraph contains only the block (and whitespace)
	paraContent := xmlStr[paraOpenClose+1 : paraCloseIdx]
	paraContent = strings.TrimSpace(paraContent)

	// If paragraph only contains the block, remove the whole paragraph
	blockContent := xmlStr[blockStart:blockEnd]
	if strings.TrimSpace(paraContent) == strings.TrimSpace(blockContent) {
		return xmlStr[:paraOpenIdx] + xmlStr[paraCloseIdx+7:]
	}

	return xmlStr
}

// performSimpleSubstitutions replaces {{key}} and {{nested.key}} patterns
func performSimpleSubstitutions(xmlStr string, ctx Context) string {
	// Find and replace simple tokens {{key}} or {{nested.key}}
	tokenRegex := regexp.MustCompile(`\{\{(\w+(?:\.\w+)*)\}\}`)

	return tokenRegex.ReplaceAllStringFunc(xmlStr, func(token string) string {
		// Extract key from {{key}}
		key := token[2 : len(token)-2]

		// Resolve value
		value := resolvePath(ctx, key)

		// Convert to string
		if value == nil {
			return ""
		}

		// Format the value
		return fmt.Sprintf("%v", value)
	})
}

// resolvePath resolves a dot-separated key path in the context
func resolvePath(ctx Context, path string) interface{} {
	parts := strings.Split(path, ".")

	var current interface{} = ctx

	for _, part := range parts {
		switch c := current.(type) {
		case Context:
			if v, exists := c[part]; exists {
				current = v
			} else {
				return nil
			}
		case map[string]interface{}:
			if v, exists := c[part]; exists {
				current = v
			} else {
				return nil
			}
		default:
			return nil
		}
	}

	return current
}

// mergeContexts creates a new context that combines parent and item fields
func mergeContexts(parent Context, item map[string]interface{}) Context {
	merged := make(Context)

	// Copy parent
	for k, v := range parent {
		merged[k] = v
	}

	// Override with item
	for k, v := range item {
		merged[k] = v
	}

	return merged
}
