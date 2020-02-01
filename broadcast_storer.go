package camforchat

var (
	CtxKeyBroadcastStorer = ContextKey("BroadcastStorer")
)

type broadcastStorer interface {
	Save(*Broadcast) error
	Find(ID string) (*Broadcast, error)
	FindByState(s BroadcastState) ([]*Broadcast, error)
	UpdateState(ID string, s BroadcastState) error
}

// GetBroadcastStorer returns broadcast storer from context
func GetBroadcastStorer(ctx context.Context) (broadcastStorer, bool) {
	l, ok := ctx.Value(CtxKeyBroadcastStorer).(broadcastStorer)
	return l, ok
}
