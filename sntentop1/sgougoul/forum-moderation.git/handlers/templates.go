package handlers

import "html/template"

// Templates is the parsed template set used by all handlers.
var Templates *template.Template

// InitTemplates parses all templates once at startup.
// This function also registers helper functions that can be used inside templates.
// Call this from main.go during application initialization.
func InitTemplates() error {

	// Register template helper functions before parsing templates.
	funcMap := template.FuncMap{
		// add allows simple integer addition inside templates.
		// Example usage in templates:
		//   {{ add $index 1 }}
		// This is used for displaying ranks like "#1, #2, #3" in lists.
		"add": func(a, b int) int {
			return a + b
		},
	}

	// Create template root, attach helper functions, then parse all templates.
	t, err := template.New("").Funcs(funcMap).ParseGlob("templates/*.html")
	if err != nil {
		return err
	}

	Templates = t
	return nil
}