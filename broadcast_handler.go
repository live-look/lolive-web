package camforchat

import (
	"context"
	"github.com/go-chi/chi/middleware"
	"log"
	"net/http"
)

var (
	ctxKeyBroadcastHandler = ContextKey("BroadcastHandler")
)

// BroadcastHandler handles by broadcasts
// keeps all online broadcast in map[string]*Broadcast (uuid => *Broadcast)
// runs new created broadcasts
type BroadcastHandler struct {
	broadcasts map[string]*Broadcast

	Publish   chan *Broadcast
	Subscribe chan *Viewer

	StopPublish   chan string
	StopSubscribe chan int64
}

// NewBroadcastHandler creates object of BroadcastHandler
func NewBroadcastHandler() *BroadcastHandler {
	return &BroadcastHandler{
		broadcasts:    make(map[string]*Broadcast),
		Publish:       make(chan *Broadcast),
		Subscribe:     make(chan *Viewer),
		StopPublish:   make(chan string),
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

				broadcast.Run(ctx)
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
func (bh *BroadcastHandler) StopBroadcasting(broadcastID string) {
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

// GetBroadcastHandler returns BroadcastHandler from context
func GetBroadcastHandler(ctx context.Context) (*BroadcastHandler, bool) {
	u, ok := ctx.Value(ctxKeyBroadcastHandler).(*BroadcastHandler)

	return u, ok
}

// BroadcastHandlerMiddleware is middleware for passing BroadcastHandler between requests
func BroadcastHandlerMiddleware(bh *BroadcastHandler) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		h := func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			r = r.WithContext(context.WithValue(r.Context(), ctxKeyBroadcastHandler, bh))
			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(h)
	}
}
