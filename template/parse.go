package template

import (
	"fmt"
	"regexp"
)

type Template struct {
	source  string
	matches map[string]string
}

// From the Python string.Template class
// https://github.com/python/cpython/blob/37c6db60f9/Lib/string.py#L60-L85
var parser = regexp.MustCompile(`(?im)\$(?:(?P<escaped>\$)|(?P<named>([_a-z][_a-z0-9]*))|{(?P<braced>([_a-z][_a-z0-9]*))}|(?P<invalid>))`)

func get(source string, start, end int) string {
	if start == -1 || end == -1 {
		return ""
	} else {
		return source[start:end]
	}
}

// New validates and creates a new template
func New(template string) (*Template, error) {
	matches := make(map[string]string)

	// Ensure all the variables are valid
	for _, submatches := range parser.FindAllStringSubmatchIndex(template, -1) {
		token := get(template, submatches[0], submatches[1])
		escaped := get(template, submatches[2], submatches[3])
		named := get(template, submatches[4], submatches[5])
		braced := get(template, submatches[8], submatches[9])

		if escaped == "$" {
			matches[token] = "$"
		} else if named != "" {
			matches[token] = named
		} else if braced != "" {
			matches[token] = braced
		} else {
			return nil, fmt.Errorf("invalid template variable at position %d", submatches[13])
		}
	}

	return &Template{source: template, matches: matches}, nil
}
