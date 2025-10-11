// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

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
