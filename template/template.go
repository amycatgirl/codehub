package templates

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
)

func safeHTML(html string) template.HTML {
	// TODO: Properly sanitize lol
	return template.HTML(html)
}

// From: https://cs.opensource.google/go/go/+/refs/tags/go1.25.5:src/html/template/template.go
func readFile(file string) (name string, b []byte, err error) {
	name = filepath.Base(file)
	b, err = os.ReadFile(file)
	return
}

func RegisterTemplatesWithSanitizer(glob string) (*template.Template, error) {
	var t *template.Template
	files, err := filepath.Glob(glob)

	if err != nil {
		return nil, err
	}

	for _, path := range files {
		name, b, err := readFile(path)
		if err != nil {
			return nil, err
		}

		s := string(b)

		var tmpl *template.Template
		if t == nil {
			t = template.New(name)
		}
		if name == t.Name() {
			tmpl = t
		} else {
			tmpl = t.New(name)
		}

		tmpl = tmpl.Funcs(template.FuncMap{
			"htmlSafe": safeHTML,
		})

		if _, err := tmpl.Parse(s); err != nil {
			return nil, fmt.Errorf("failed to parse template %s: %w", path, err)
		}
	}

	return t, nil
}
