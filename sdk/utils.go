package sdk

import (
	"regexp"
)

var nextLinkRE = regexp.MustCompile(`<([^>]+)>;\s*rel="next"`)

// extractNextLink parses the HTTP Link header and returns the "next" URL, if present.
func extractNextLink(linkHeader string) (string, error) {
	if linkHeader == "" {
		return "", nil // No pagination
	}

	matches := nextLinkRE.FindStringSubmatch(linkHeader)
	if len(matches) != 2 {
		return "", nil // No next link found; not an error
	}

	return matches[1], nil
}
