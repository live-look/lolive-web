package camforchat

import (
	"context"
	"github.com/go-chi/chi/middleware"
	"github.com/pion/webrtc/v2"
	"net/http"
)

var (
	ctxKeyWebrtcAPI = ContextKey("WebrtcAPI")

	webrtcConfig = webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}
)

// Webrtc is global webrtc api
type Webrtc struct {
	API *webrtc.API
}

// NewWebrtc creates webrtc object
func NewWebrtc() *Webrtc {
	m := webrtc.MediaEngine{}
	m.RegisterCodec(webrtc.NewRTPVP8Codec(webrtc.DefaultPayloadTypeVP8, 90000))
	// m.RegisterCodec(webrtc.NewRTPOpusCodec(webrtc.DefaultPayloadTypeOpus, 48000))

	return &Webrtc{
		API: webrtc.NewAPI(webrtc.WithMediaEngine(m)),
	}
}

// NewPeerConnection creates new peerconnection
func (w *Webrtc) NewPeerConnection() (*webrtc.PeerConnection, error) {
	return w.API.NewPeerConnection(webrtcConfig)
}

// GetWebrtcAPI returns database connection link
func GetWebrtcAPI(ctx context.Context) (*Webrtc, bool) {
	w, ok := ctx.Value(ctxKeyWebrtcAPI).(*Webrtc)
	return w, ok
}

// WebrtcAPIMiddleware is middleware for passing Webrtc between requests
func WebrtcAPIMiddleware(wrtc *Webrtc) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		h := func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			r = r.WithContext(context.WithValue(r.Context(), ctxKeyWebrtcAPI, wrtc))
			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(h)
	}
}
