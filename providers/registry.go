package providers

import (
	"strings"
)

type creator func(id string) (Provider, error)

var registry = map[string]creator{"mailgun": NewMailgun}

// Get retrieves and creates the provider
func Get(id, typeName string) (Provider, error) {
	providerCreator, ok := registry[strings.ToLower(typeName)]
	if !ok {
		return nil, nil
	}

	return providerCreator(id)
}
