package models

import (
	"github.com/pion/webrtc/v2"
)

// Webrtc is global webrtc api
type Webrtc struct {
	API *webrtc.API
}

// NewWebrtc creates webrtc object
func NewWebrtc() *Webrtc {
	m := webrtc.MediaEngine{}
	m.RegisterCodec(webrtc.NewRTPVP8Codec(webrtc.DefaultPayloadTypeVP8, 90000))

	return &Webrtc{API: webrtc.NewAPI(webrtc.WithMediaEngine(m))}
}
