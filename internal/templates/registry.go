package templates

import (
	"sort"
	"strings"
)

// Registry holds the built-in templates.
var Registry = map[string]string{
	"minimal-go": `
project-root/
├── cmd/
│   └── main.go
├── internal/
├── go.mod
└── README.md
`,
	"react-vite": `
my-app/
├── public/
│   └── vite.svg
├── src/
│   ├── assets/
│   ├── components/
│   ├── App.css
│   ├── App.tsx
│   ├── index.css
│   └── main.tsx
├── index.html
├── package.json
├── tsconfig.json
└── vite.config.ts
`,
	"python-flask": `
flask-app/
├── app/
│   ├── templates/
│   │   └── index.html
│   ├── static/
│   │   └── style.css
│   ├── __init__.py
│   └── routes.py
├── tests/
├── venv/
├── config.py
├── requirements.txt
└── run.py
`,
}

// List returns a sorted list of available template names.
func List() []string {
	keys := make([]string, 0, len(Registry))
	for k := range Registry {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// Get returns the content of a template by name.
func Get(name string) (string, bool) {
	content, ok := Registry[name]
	return strings.TrimSpace(content), ok
}
