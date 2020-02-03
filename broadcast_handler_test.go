package camforchat

import (
	"context"
	"testing"
	"time"
)

func TestStopView(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	bh := NewBroadcastHandler()
	viewerIDs := make(chan int64)

	go func() {
		ID := <-bh.StopSubscribe
		viewerIDs <- ID
	}()

	bh.StopView(1)

	select {
	case viewerID := <-viewerIDs:
		if viewerID != 1 {
			t.Fatalf("Expected viewerID is eq 1, but not: %v", viewerID)
		}
	case <-ctx.Done():
		t.Fatalf("Timeout: %v", ctx.Err())
	}
}
