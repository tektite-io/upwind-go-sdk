// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package sdk

import (
	"context"
)

// CollectAll is a helper function that collects all items from a channel into a slice
func CollectAll[T any](ctx context.Context, itemsCh <-chan T, errCh <-chan error) ([]T, error) {
	var results []T

	for {
		select {
		case item, ok := <-itemsCh:
			if !ok {
				// Channel closed, check for errors
				select {
				case err := <-errCh:
					if err != nil {
						return results, err
					}
				default:
				}
				return results, nil
			}
			results = append(results, item)
		case err := <-errCh:
			if err != nil {
				return results, err
			}
		case <-ctx.Done():
			return results, ctx.Err()
		}
	}
}
