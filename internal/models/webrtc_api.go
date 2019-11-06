package models

import (
	"github.com/pion/webrtc/v2"
)

// Webrtc is global webrtc api
type Webrtc struct {
	API *webrtc.API
}

var (
	webrtcConfig = webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}
)

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
