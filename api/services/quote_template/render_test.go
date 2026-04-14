package quote_template

import (
	"strings"
	"testing"
)

func TestSimpleSubstitution(t *testing.T) {
	ctx := Context{
		"name": "John Doe",
		"age":  "30",
	}

	xml := `<w:p><w:r><w:t>Hello {{name}}, you are {{age}} years old</w:t></w:r></w:p>`
	result := performSimpleSubstitutions(xml, ctx)

	if !strings.Contains(result, "Hello John Doe") {
		t.Errorf("Name not substituted: %s", result)
	}
	if !strings.Contains(result, "30 years old") {
		t.Errorf("Age not substituted: %s", result)
	}
	if strings.Contains(result, "{{") {
		t.Errorf("Template tokens still present: %s", result)
	}
}

func TestNestedPathSubstitution(t *testing.T) {
	ctx := Context{
		"insurer": map[string]interface{}{
			"name":  "Insurance Co",
			"email": "contact@insurer.com",
		},
	}

	xml := `<w:p><w:r><w:t>{{insurer.name}} - {{insurer.email}}</w:t></w:r></w:p>`
	result := performSimpleSubstitutions(xml, ctx)

	if !strings.Contains(result, "Insurance Co") {
		t.Errorf("Nested name not substituted: %s", result)
	}
	if !strings.Contains(result, "contact@insurer.com") {
		t.Errorf("Nested email not substituted: %s", result)
	}
}

func TestMissingKeyHandling(t *testing.T) {
	ctx := Context{
		"name": "John",
	}

	xml := `<w:p><w:r><w:t>{{name}} {{missing}}</w:t></w:r></w:p>`
	result := performSimpleSubstitutions(xml, ctx)

	if !strings.Contains(result, "John") {
		t.Errorf("Name not substituted: %s", result)
	}
	// Missing key should be replaced with empty string
	if strings.Contains(result, "{{missing}}") {
		t.Errorf("Missing key not handled: %s", result)
	}
}

func TestConditionalBlockTrue(t *testing.T) {
	ctx := Context{
		"show_details": true,
	}

	xml := `<w:p><w:r><w:t>Start {{#show_details}}Details here{{/show_details}} End</w:t></w:r></w:p>`
	result := processBlocks(xml, ctx)

	if !strings.Contains(result, "Details here") {
		t.Errorf("Block not included when true: %s", result)
	}
	if strings.Contains(result, "{{#") {
		t.Errorf("Block markers not removed: %s", result)
	}
}

func TestConditionalBlockFalse(t *testing.T) {
	ctx := Context{
		"show_details": false,
	}

	xml := `<w:p><w:r><w:t>Start {{#show_details}}Details here{{/show_details}} End</w:t></w:r></w:p>`
	result := processBlocks(xml, ctx)

	if strings.Contains(result, "Details here") {
		t.Errorf("Block should be excluded when false: %s", result)
	}
}

func TestIterationBlock(t *testing.T) {
	ctx := Context{
		"items": []map[string]interface{}{
			{"name": "Item1", "value": "100"},
			{"name": "Item2", "value": "200"},
		},
	}

	xml := `<w:p><w:r><w:t>{{#items}}{{name}}: {{value}}, {{/items}}</w:t></w:r></w:p>`
	result := processBlocks(xml, ctx)

	if !strings.Contains(result, "Item1") {
		t.Errorf("First item not found: %s", result)
	}
	if !strings.Contains(result, "Item2") {
		t.Errorf("Second item not found: %s", result)
	}
	if !strings.Contains(result, "100") {
		t.Errorf("First value not found: %s", result)
	}
	if !strings.Contains(result, "200") {
		t.Errorf("Second value not found: %s", result)
	}
}

func TestResolvePath(t *testing.T) {
	tests := []struct {
		name     string
		ctx      Context
		path     string
		expected interface{}
	}{
		{
			name:     "Simple key",
			ctx:      Context{"name": "John"},
			path:     "name",
			expected: "John",
		},
		{
			name: "Nested key",
			ctx: Context{
				"user": map[string]interface{}{
					"name": "Jane",
				},
			},
			path:     "user.name",
			expected: "Jane",
		},
		{
			name:     "Missing key",
			ctx:      Context{"name": "John"},
			path:     "missing",
			expected: nil,
		},
		{
			name: "Deep nesting",
			ctx: Context{
				"level1": map[string]interface{}{
					"level2": map[string]interface{}{
						"value": "deep",
					},
				},
			},
			path:     "level1.level2.value",
			expected: "deep",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := resolvePath(tt.ctx, tt.path)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestMergeContexts(t *testing.T) {
	parent := Context{
		"name":   "Parent",
		"age":    "40",
		"nested": map[string]interface{}{"value": "parent"},
	}

	item := map[string]interface{}{
		"name": "Child",
		"id":   "123",
	}

	merged := mergeContexts(parent, item)

	// Item values should override parent
	if merged["name"] != "Child" {
		t.Errorf("Item did not override parent: %v", merged["name"])
	}

	// Parent values should remain if not overridden
	if merged["age"] != "40" {
		t.Errorf("Parent value lost: %v", merged["age"])
	}

	// New item keys should be present
	if merged["id"] != "123" {
		t.Errorf("New item key missing: %v", merged["id"])
	}
}
