// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package sdk

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestCollectAll(t *testing.T) {
	t.Run("successful collection", func(t *testing.T) {
		ctx := context.Background()
		itemsCh := make(chan int, 5)
		errCh := make(chan error, 1)

		// Send items
		go func() {
			for i := 1; i <= 5; i++ {
				itemsCh <- i
			}
			close(itemsCh)
			close(errCh)
		}()

		results, err := CollectAll(ctx, itemsCh, errCh)
		if err != nil {
			t.Fatalf("CollectAll() error = %v", err)
		}

		if len(results) != 5 {
			t.Errorf("CollectAll() got %d items, want 5", len(results))
		}

		for i, val := range results {
			if val != i+1 {
				t.Errorf("CollectAll() item %d = %d, want %d", i, val, i+1)
			}
		}
	})

	t.Run("with error", func(t *testing.T) {
		ctx := context.Background()
		itemsCh := make(chan int, 5)
		errCh := make(chan error, 1)

		expectedErr := errors.New("test error")

		// Send items and error
		go func() {
			for i := 1; i <= 3; i++ {
				itemsCh <- i
			}
			errCh <- expectedErr
			close(itemsCh)
		}()

		results, err := CollectAll(ctx, itemsCh, errCh)
		if err != expectedErr {
			t.Errorf("CollectAll() error = %v, want %v", err, expectedErr)
		}

		// Should have received some items before error
		if len(results) < 1 {
			t.Errorf("CollectAll() got %d items, want at least 1", len(results))
		}
	})

	t.Run("context cancellation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		itemsCh := make(chan int)
		errCh := make(chan error, 1)

		// Cancel immediately
		cancel()

		go func() {
			time.Sleep(100 * time.Millisecond)
			close(itemsCh)
			close(errCh)
		}()

		_, err := CollectAll(ctx, itemsCh, errCh)
		if err != context.Canceled {
			t.Errorf("CollectAll() error = %v, want %v", err, context.Canceled)
		}
	})

	t.Run("empty channel", func(t *testing.T) {
		ctx := context.Background()
		itemsCh := make(chan int)
		errCh := make(chan error, 1)

		// Close immediately
		close(itemsCh)
		close(errCh)

		results, err := CollectAll(ctx, itemsCh, errCh)
		if err != nil {
			t.Fatalf("CollectAll() error = %v", err)
		}

		if len(results) != 0 {
			t.Errorf("CollectAll() got %d items, want 0", len(results))
		}
	})
}
