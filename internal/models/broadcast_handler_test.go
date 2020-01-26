package models

import (
	"testing"
)

func TestStopView(t *testing.T) {
	bh := NewBroadcastHandler()

	var viewerID int64

	go func() {
		viewerID = <-bh.StopSubscribe
	}()

	bh.StopView(1)

	if viewerID != 1 {
		t.Fatalf("Expected viewerID is eq 1, but not: %v", viewerID)
	}
}
