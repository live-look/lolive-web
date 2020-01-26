package models

import (
	"camforchat/internal/usecases"
	"context"
	"github.com/google/uuid"
	"github.com/pion/rtcp"
	"github.com/pion/webrtc/v2"
	"io"
	"log"
	"time"
)

// BroadcastState is global state of broadcast
type BroadcastState string

const (
	// BroadcastStateOnline when broadcaster is online
	BroadcastStateOnline BroadcastState = "online"
	// BroadcastStateOffline when broadcaster is off
	BroadcastStateOffline BroadcastState = "offline"
	// BroadcastStatePrivate when broadcaster is in private
	BroadcastStatePrivate BroadcastState = "private"
)

// Broadcast is struct of broadcaster
type Broadcast struct {
	ID        string         `json:"id" db:"id"`
	UserID    int64          `json:"user_id" db:"user_id"`
	CreatedAt time.Time      `json:"created_at" db:"created_at"`
	State     BroadcastState `json:"state" db:"state"`

	UserName string `json:"user_name" db:"user_name"`

	// SDP from client
	LocalSessionDescription string `json:"local_sdp" db:"-"`
	// SDP on server
	RemoteSessionDescription string `json:"remote_sdp" db:"-"`

	SDPChan chan string  `json:"-" db:"-"`
	Publish chan *Viewer `json:"-" db:"-"`

	viewers map[int64]*Viewer

	webrtc *Webrtc

	// Connection with browser (webcam of the broadcaster)
	peerConnection *webrtc.PeerConnection

	screenshotTaker *BroadcastScreenshot
}

// NewBroadcast creates new instance of models.Broadcast
// user - instance of User model
// webrtc - instance confgured Webrtc API
func NewBroadcast(user *User, webrtc *webrtc) (*Broadcast, error) {
	var err error

	bc := &Broadcast{
		ID:    uuid.New().String(),
		State: BroadcastStateOffline,

		UserID:    user.ID,
		UserName:  user.Name,
		CreatedAt: time.Now(),

		SDPChan: make(chan string),
		Publish: make(chan *Viewer),

		viewers: make(map[int64]*Viewer),
		webrtc:  webrtc,
	}

	bc.screenshotTaker, err = NewBroadcastScreenshot(bc.ID)
	if err != nil {
		return nil, err
	}

	bc.peerConnection, err = webrtc.NewPeerConnection()
	if err != nil {
		return nil, err
	}

	return bc, nil
}

// Join joins viewer to broadcast
func (b *Broadcast) Join(viewer *Viewer) {
	b.Publish <- viewer
}

// Run starts broadcast loop
func (b *Broadcast) Run() {
	var err error

	// Run screenshot taker loop
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	b.screenshotTaker.TakeScreenshots(ctx)

	// Run transmission loop
	go func(broadcast *Broadcast) {
		log.Println("Start broadcasting...")

		if _, err = broadcast.peerConnection.AddTransceiver(webrtc.RTPCodecTypeVideo); err != nil {
			log.Println(err)
			return
		}

		offer := webrtc.SessionDescription{}
		usecases.DecodeSDP(broadcast.LocalSessionDescription, &offer)

		// Set the remote SessionDescription
		err = peerConnection.SetRemoteDescription(offer)
		if err != nil {
			log.Println(err)
			return
		}

		// Create answer
		answer, err := peerConnection.CreateAnswer(nil)
		if err != nil {
			log.Println(err)
			return
		}

		// Sets the LocalDescription, and starts our UDP listeners
		err = peerConnection.SetLocalDescription(answer)
		if err != nil {
			log.Println(err)
			return
		}

		peerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
			log.Printf("Connection State has changed %s \n", connectionState.String())

			if connectionState == webrtc.ICEConnectionStateConnected {
				log.Println("Connected")
				err := broadcast.SetState(BroadcastStateOnline)
				if err != nil {
					log.Printf("Error set state: %v\n", err)
				}
			} else if connectionState == webrtc.ICEConnectionStateFailed ||
				connectionState == webrtc.ICEConnectionStateDisconnected {
				log.Println("Disconnected")
				err := broadcast.SetState(BroadcastStateOffline)
				if err != nil {
					log.Printf("Error set state: %v\n", err)
				}
			}
		})

		localTrackChan := make(chan *webrtc.Track)

		peerConnection.OnTrack(func(track *webrtc.Track, receiver *webrtc.RTPReceiver) {
			go func() {
				ticker := time.NewTicker(time.Second * 3)
				for range ticker.C {
					errSend := peerConnection.WriteRTCP([]rtcp.Packet{&rtcp.PictureLossIndication{MediaSSRC: track.SSRC()}})
					if errSend != nil {
						log.Println(errSend)
					}
				}
			}()

			log.Println("On track")

			localTrack, newTrackErr := peerConnection.NewTrack(track.PayloadType(), track.SSRC(), "video", "pion")
			if newTrackErr != nil {
				log.Printf("Error: %v", newTrackErr)
				return
			}
			localTrackChan <- localTrack

			rtpBuf := make([]byte, 1200)

			for {
				i, readErr := track.Read(rtpBuf)
				if readErr != nil {
					log.Printf("Error: %v", readErr)
					break
				}
				buf := rtpBuf[:i]

				// ErrClosedPipe means we don't have any subscribers, this is ok if no peers have connected yet
				if _, err = localTrack.Write(buf); err != nil && err != io.ErrClosedPipe {
					log.Printf("Error: %v", err)
					break
				}

				if track.Codec().Name == webrtc.VP8 {
					broadcast.screenshotTaker.StreamBuf <- buf
				}
			}
		})

		broadcast.SDPChan <- usecases.EncodeSDP(answer)
		localTrack := <-localTrackChan

		for {
			select {
			case viewer := <-b.Publish:
				log.Printf("Viewer %d joined to broadcast", viewer.ID)
				b.viewers[viewer.ID] = viewer
				viewer.Run(localTrack)
			}
		}
	}(b)
}

func (b *Broadcast) Stop() chan (bool) {

}
