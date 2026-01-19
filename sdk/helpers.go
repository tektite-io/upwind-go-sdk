// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package sdk

import (
	"context"
	"runtime"
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

// CollectInChunks collects items from a channel into chunks for memory-efficient processing
// It yields chunks of the specified size via a callback function
// This is useful for processing large datasets without loading everything into memory
//
// Example:
//
//	findingsCh, errCh := client.ListVulnerabilityFindings(ctx, query)
//	err := sdk.CollectInChunks(ctx, findingsCh, errCh, 1000, func(chunk []VulnerabilityFinding) error {
//	    // Process chunk (e.g., write to database)
//	    return processChunk(chunk)
//	})
func CollectInChunks[T any](ctx context.Context, itemsCh <-chan T, errCh <-chan error, chunkSize int, processChunk func([]T) error) error {
	chunk := make([]T, 0, chunkSize)
	chunkCount := 0

	for {
		select {
		case item, ok := <-itemsCh:
			if !ok {
				// Channel closed, process remaining items
				if len(chunk) > 0 {
					if err := processChunk(chunk); err != nil {
						return err
					}
				}
				// Check for errors
				select {
				case err := <-errCh:
					return err
				default:
					return nil
				}
			}

			chunk = append(chunk, item)

			// Process chunk when it reaches target size
			if len(chunk) >= chunkSize {
				if err := processChunk(chunk); err != nil {
					return err
				}
				chunkCount++

				// Clear chunk and release memory
				chunk = make([]T, 0, chunkSize)

				// Trigger GC every 10 chunks to manage memory
				if chunkCount%10 == 0 {
					runtime.GC()
				}
			}
		case err := <-errCh:
			if err != nil {
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// StreamInChunks streams items from a channel into chunks via a new channel
// This allows for composable, memory-efficient processing pipelines
//
// Example:
//
//	findingsCh, errCh := client.ListVulnerabilityFindings(ctx, query)
//	chunksCh := sdk.StreamInChunks(ctx, findingsCh, 1000)
//	for chunk := range chunksCh {
//	    processChunk(chunk)
//	}
//	if err := <-errCh; err != nil {
//	    log.Fatal(err)
//	}
func StreamInChunks[T any](ctx context.Context, itemsCh <-chan T, chunkSize int) <-chan []T {
	chunksCh := make(chan []T, 10)

	go func() {
		defer close(chunksCh)

		chunk := make([]T, 0, chunkSize)
		chunkCount := 0

		for {
			select {
			case item, ok := <-itemsCh:
				if !ok {
					// Channel closed, send remaining items
					if len(chunk) > 0 {
						select {
						case chunksCh <- chunk:
						case <-ctx.Done():
							return
						}
					}
					return
				}

				chunk = append(chunk, item)

				// Send chunk when it reaches target size
				if len(chunk) >= chunkSize {
					select {
					case chunksCh <- chunk:
						chunkCount++
						chunk = make([]T, 0, chunkSize)

						// Trigger GC every 10 chunks
						if chunkCount%10 == 0 {
							runtime.GC()
						}
					case <-ctx.Done():
						return
					}
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	return chunksCh
}
