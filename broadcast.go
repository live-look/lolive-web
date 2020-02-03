package camforchat

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/pion/rtcp"
	"github.com/pion/webrtc/v2"
	"gitlab.com/isqad/camforchat/utils"
	"go.uber.org/zap"
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

	// Connection with browser (webcam of the broadcaster)
	peerConnection *webrtc.PeerConnection

	screenshotTaker *BroadcastScreenshot
}

// NewBroadcast creates new instance of models.Broadcast
// user - instance of User model
func NewBroadcast(userID int64, userName string, webrtcAPI *Webrtc) (*Broadcast, error) {
	var err error

	bc := &Broadcast{
		ID:    uuid.New().String(),
		State: BroadcastStateOffline,

		UserID:    userID,
		UserName:  userName,
		CreatedAt: time.Now(),

		SDPChan: make(chan string),
		Publish: make(chan *Viewer),

		viewers: make(map[int64]*Viewer),
	}

	bc.screenshotTaker, err = NewBroadcastScreenshot(bc.ID)
	if err != nil {
		return nil, err
	}

	bc.peerConnection, err = webrtcAPI.NewPeerConnection()
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
// context must keep logger, BroadcastStorer
func (b *Broadcast) Run(ctx context.Context) error {
	var err error

	logger, ok := GetLogger(ctx)
	if !ok {
		return errors.New("No logger in context")
	}

	storer, ok := GetBroadcastStorer(ctx)
	if !ok {
		return errors.New("No broadcast storer in context")
	}

	// Run screenshot taker loop
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	b.screenshotTaker.TakeScreenshots(ctx)

	// Run transmission loop
	go func(broadcast *Broadcast) {
		logger.Info("Starting broadcast", zap.String("ID", broadcast.ID))

		if _, err = broadcast.peerConnection.AddTransceiver(webrtc.RTPCodecTypeVideo); err != nil {
			logger.Error("Failed add transceiver", zap.String("ID", broadcast.ID), zap.Error(err))

			return
		}

		offer := webrtc.SessionDescription{}
		utils.DecodeSDP(broadcast.LocalSessionDescription, &offer)

		// Set the remote SessionDescription
		err = broadcast.peerConnection.SetRemoteDescription(offer)
		if err != nil {
			logger.Error("Failed set remote description", zap.String("ID", broadcast.ID), zap.Error(err))
			return
		}

		// Create answer
		answer, err := broadcast.peerConnection.CreateAnswer(nil)
		if err != nil {
			logger.Error("Failed create answer", zap.String("ID", broadcast.ID), zap.Error(err))
			return
		}

		// Sets the LocalDescription, and starts our UDP listeners
		err = broadcast.peerConnection.SetLocalDescription(answer)
		if err != nil {
			logger.Error("Failed set local description", zap.String("ID", broadcast.ID), zap.Error(err))
			return
		}

		broadcast.peerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
			logger.Info("Connection state has changed", zap.String("ID", broadcast.ID), zap.String("state", connectionState.String()))

			if connectionState == webrtc.ICEConnectionStateConnected {
				logger.Info("Broadcast is online", zap.String("ID", broadcast.ID))

				err := storer.UpdateState(broadcast.ID, BroadcastStateOnline)
				if err != nil {
					logger.Error("Failed update state", zap.String("ID", broadcast.ID), zap.Error(err))
				}
			} else if connectionState == webrtc.ICEConnectionStateFailed ||
				connectionState == webrtc.ICEConnectionStateDisconnected {

				logger.Info("Broadcast is off", zap.String("ID", broadcast.ID))

				err := storer.UpdateState(broadcast.ID, BroadcastStateOffline)
				if err != nil {
					logger.Error("Failed update state", zap.String("ID", broadcast.ID), zap.Error(err))
				}
			}
		})

		localTrackChan := make(chan *webrtc.Track)

		broadcast.peerConnection.OnTrack(func(track *webrtc.Track, receiver *webrtc.RTPReceiver) {
			go func() {
				ticker := time.NewTicker(time.Second * 3)
				for range ticker.C {
					errSend := broadcast.peerConnection.WriteRTCP([]rtcp.Packet{&rtcp.PictureLossIndication{MediaSSRC: track.SSRC()}})
					if errSend != nil {
						log.Println(errSend)
					}
				}
			}()

			logger.Info("Broadcast on track", zap.String("ID", broadcast.ID))

			localTrack, newTrackErr := broadcast.peerConnection.NewTrack(track.PayloadType(), track.SSRC(), "video", "pion")
			if newTrackErr != nil {
				logger.Error("Failed create new track", zap.String("ID", broadcast.ID), zap.Error(newTrackErr))
				return
			}
			localTrackChan <- localTrack

			rtpBuf := make([]byte, 1200)

			for {
				i, readErr := track.Read(rtpBuf)
				if readErr != nil {
					logger.Error("Failed reading from track", zap.String("ID", broadcast.ID), zap.Error(readErr))
					break
				}
				buf := rtpBuf[:i]

				// ErrClosedPipe means we don't have any subscribers, this is ok if no peers have connected yet
				if _, err = localTrack.Write(buf); err != nil && err != io.ErrClosedPipe {
					logger.Error("Failed writing track", zap.String("ID", broadcast.ID), zap.Error(readErr))
					break
				}

				if track.Codec().Name == webrtc.VP8 {
					broadcast.screenshotTaker.StreamBuf <- buf
				}
			}
		})

		broadcast.SDPChan <- utils.EncodeSDP(answer)
		localTrack := <-localTrackChan

		for {
			select {
			case viewer := <-b.Publish:
				logger.Info("Viewer joined to broadcast", zap.String("ID", broadcast.ID), zap.Int64("Viewer ID", viewer.ID))

				b.viewers[viewer.ID] = viewer
				viewer.Run(localTrack)
			}
		}
	}(b)

	return nil
}

// Stop stops broadcasting
func (b *Broadcast) Stop() chan (bool) {
	return make(chan bool)
}
