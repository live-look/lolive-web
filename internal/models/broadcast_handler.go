package models

import (
	"context"
	"log"
)

// BroadcastHandler handles by broadcasts
// keeps all online broadcast in map[int64]*Broadcast
// runs new created broadcasts
type BroadcastHandler struct {
	broadcasts map[int64]*Broadcast

	Publish   chan *Broadcast
	Subscribe chan *Viewer
}

// NewBroadcastHandler creates object of BroadcastHandler
func NewBroadcastHandler() *BroadcastHandler {
	return &BroadcastHandler{
		broadcasts: make(map[int64]*Broadcast),
		Publish:    make(chan *Broadcast),
	}
}

// Run starts main loop of handler
func (bh *BroadcastHandler) Run(ctx context.Context) {
	go func(bh *BroadcastHandler, ctx context.Context) {
		log.Println("Starting broadcasting handler...")

		for {
			select {
			case broadcast := <-bh.Publish:
				bh.broadcasts[broadcast.ID] = broadcast

				broadcast.Run()
			case viewer := <-bh.Subscribe:
				log.Println("Found broadcast")

				broadcast := bh.broadcasts[viewer.BroadcastID]
				broadcast.Join(viewer)
			case <-ctx.Done():
				// TODO: graceful shutdown all broadcasts
				break
			}
		}
	}(bh, ctx)
}

// StartBroadcasting sends created broadcast to Publish channel for add to broadcast register
// and run its main loop
func (bh *BroadcastHandler) StartBroadcasting(broadcast *Broadcast) {
	bh.Publish <- broadcast
}

// StartView subsrcibes new viewer to given broadcast (property of viewer)
func (bh *BroadcastHandler) StartView(viewer *Viewer) {
	bh.Subscribe <- viewer
}
