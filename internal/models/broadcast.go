package models

import (
	"camforchat/internal/usecases"
	"context"
	"github.com/jmoiron/sqlx"
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
	ID        int64          `json:"id" db:"id"`
	UserID    int64          `json:"user_id" db:"user_id"`
	CreatedAt time.Time      `json:"created_at" db:"created_at"`
	State     BroadcastState `json:"state" db:"state"`

	UserName string `json:"user_name" db:"user_name"`

	screenshotTaker *BroadcastScreenshot `json:"-" db:"-"`

	LocalSessionDescription  string `json:"local_sdp" db:"-"`
	RemoteSessionDescription string `json:"remote_sdp" db:"-"`

	SDPChan chan string  `json:"-" db:"-"`
	Publish chan *Viewer `json:"-" db:"-"`

	db *sqlx.DB

	viewers map[int64]*Viewer

	webrtc *Webrtc
}

// NewBroadcast creates new instance of models.broadcast
func NewBroadcast(db *sqlx.DB, webrtc *Webrtc, userID int64) *Broadcast {
	return &Broadcast{
		db:      db,
		UserID:  userID,
		SDPChan: make(chan string),
		Publish: make(chan *Viewer),
		viewers: make(map[int64]*Viewer),
		State:   BroadcastStateOffline,
		webrtc:  webrtc,
	}
}

// FindBroadcast gets broadcast from db by ID
func FindBroadcast(db *sqlx.DB, webrtc *Webrtc, ID int64) (*Broadcast, error) {
	broadcast := NewBroadcast(db, webrtc, 0)

	selectQuery := `SELECT b.*, u.name AS user_name
					FROM broadcasts b
					INNER JOIN users u
					  ON u.id = b.user_id
					WHERE b.id = $1`

	err := db.Get(broadcast, selectQuery, ID)

	return broadcast, err
}

// GetBroadcastsByState retrive list of online broadcasts
func GetBroadcastsByState(db *sqlx.DB, state BroadcastState) ([]*Broadcast, error) {
	var broadcasts []*Broadcast
	selectQuery := `SELECT b.*, u.name AS user_name FROM broadcasts b INNER JOIN users u ON u.id = b.user_id WHERE state = $1 ORDER BY created_at DESC`

	err := db.Select(&broadcasts, selectQuery, state)

	return broadcasts, err
}

// Save saves broadcast into the database
func (b *Broadcast) Save(broadcastUser *User) error {
	b.UserName = broadcastUser.Name
	b.CreatedAt = time.Now()

	insertQuery := `INSERT INTO broadcasts (user_id, state, created_at) VALUES ($1, $2, NOW()) RETURNING id`
	if err := b.db.Get(&b.ID, insertQuery, broadcastUser.ID, b.State); err != nil {
		return err
	}
	return nil
}

// SetState changes state of broadcast
func (b *Broadcast) SetState(state BroadcastState) error {
	updateQuery := `UPDATE broadcasts SET state = :state WHERE id = :id`
	_, err := b.db.NamedExec(updateQuery,
		map[string]interface{}{
			"state": state,
			"id":    b.ID,
		})

	return err
}

// Join joins viewer to broadcast
func (b *Broadcast) Join(viewer *Viewer) {
	b.Publish <- viewer
}

// Run starts broadcast loop
func (b *Broadcast) Run() {
	var err error
	// Run screenshot taker loop
	b.screenshotTaker, err = NewBroadcastScreenshot(b.ID)
	if err != nil {
		log.Println("Fail to run: ", err)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	b.screenshotTaker.TakeScreenshots(ctx)

	// Run transmission loop
	go func(broadcast *Broadcast) {
		log.Println("Start broadcasting...")

		peerConnection, err := broadcast.webrtc.NewPeerConnection()
		if err != nil {
			log.Println(err)
			return
		}

		if _, err = peerConnection.AddTransceiver(webrtc.RTPCodecTypeVideo); err != nil {
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
