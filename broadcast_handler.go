package camforchat

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

	StopPublish   chan int64
	StopSubscribe chan int64
}

// NewBroadcastHandler creates object of BroadcastHandler
func NewBroadcastHandler() *BroadcastHandler {
	return &BroadcastHandler{
		broadcasts:    make(map[int64]*Broadcast),
		Publish:       make(chan *Broadcast),
		Subscribe:     make(chan *Viewer),
		StopPublish:   make(chan int64),
		StopSubscribe: make(chan int64),
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
				log.Println("broadcast runned")
			case viewer := <-bh.Subscribe:
				log.Println("Found broadcast")

				broadcast := bh.broadcasts[viewer.BroadcastID]
				broadcast.Join(viewer)
			case broadcastID := <-bh.StopPublish:
				// stop broadcast and remove it
				<-bh.broadcasts[broadcastID].Stop()
				delete(bh.broadcasts, broadcastID)
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

// StopBroadcasting stops main loop of broadcaster and removes him
func (bh *BroadcastHandler) StopBroadcasting(broadcastID int64) {
	bh.StopPublish <- broadcastID
}

// StartView subsrcibes new viewer to given broadcast (property of viewer)
func (bh *BroadcastHandler) StartView(viewer *Viewer) {
	bh.Subscribe <- viewer
}

// StopView stops main loop of viewer and removes him
func (bh *BroadcastHandler) StopView(viewerID int64) {
	bh.StopSubscribe <- viewerID
}
