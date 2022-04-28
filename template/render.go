package template

import (
	"strings"
	"text/template"
)

func (t *Template) Render(context map[string]string, html bool) string {
	result := t.source

	// Replace all the tokens
	for token, key := range t.matches {
		// Get the value for the token
		var value string
		if token == "$$" {
			value = "$"
		} else {
			var ok bool
			value, ok = context[key]
			if !ok {
				continue
			}

			// Escape the value if necessary
			if html {
				value = template.HTMLEscapeString(value)
			}
		}

		result = strings.ReplaceAll(result, token, value)
	}

	return result
}
