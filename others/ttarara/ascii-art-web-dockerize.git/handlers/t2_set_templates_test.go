package handlers

import (
	"html/template"
	"testing"
)

// TestSetTemplates verifies that SetTemplates correctly assigns the global tpl variable.
func TestSetTemplates(t *testing.T) {

	originalTpl := tpl
	defer func() { tpl = originalTpl }()

	// Create a dummy template
	tmpl, err := template.New("test").Parse("<h1>{{.Title}}</h1>")
	if err != nil {
		t.Fatalf("Failed to create test template: %v", err)
	}

	// Call SetTemplates
	SetTemplates(tmpl)

	// Check if tpl was set correctly
	if tpl != tmpl {
		t.Errorf("Expected tpl to be set to %v, but got %v", tmpl, tpl)
	}
}

func TestSetTemplates_Nil(t *testing.T) {

	originalTpl := tpl
	defer func() { tpl = originalTpl }()

	SetTemplates(nil)
	if tpl != nil {
		t.Errorf("Expected tpl to be nil, but got %v", tpl)
	}
}
