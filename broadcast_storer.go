package camforchat

import (
	"context"
)

var (
	ctxKeyBroadcastStorer = ContextKey("BroadcastStorer")
)

// BroadcastStorer is interface for storing Broadcast
type BroadcastStorer interface {
	Save(*Broadcast) error
	Find(ID string) (*Broadcast, error)
	FindByState(s BroadcastState) ([]*Broadcast, error)
	UpdateState(ID string, s BroadcastState) error
}

// GetBroadcastStorer returns broadcast storer from context
func GetBroadcastStorer(ctx context.Context) (BroadcastStorer, bool) {
	l, ok := ctx.Value(ctxKeyBroadcastStorer).(BroadcastStorer)
	return l, ok
}
