package utils

import (
	"fmt"
	"net/url"
)

// ValidateOriginURLsFromArray takes an array of strings containing URLs.
// If any of the URLs can't be parsed, it returns an error.
func ValidateOriginURLsFromArray(urls []string) error {
	if len(urls) == 1 {
		if urls[0] == "*" {
			return nil
		}
	}
	for _, u := range urls {
		if _, err := url.ParseRequestURI(u); err != nil {
			return fmt.Errorf("Error parsing origin '%s': %w", u, err)
		}
	}
	return nil
}
