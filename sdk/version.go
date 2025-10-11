// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package sdk

// Version is the current version of the Upwind Go SDK (format: vX.Y.Z)
const Version = "v1.0.1"

// UserAgent returns the User-Agent string for HTTP requests
func UserAgent() string {
	return "upwind-go-sdk/" + Version
}
