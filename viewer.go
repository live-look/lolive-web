package camforchat

import (
	"github.com/jmoiron/sqlx"
	"github.com/pion/webrtc/v2"
	"gitlab.com/isqad/camforchat/utils"
	"log"
	"time"
)

// ViewerState is states of viewer
type ViewerState string

const (
	// ViewerStateJoined is mean that member has been joined to broadcast
	ViewerStateJoined ViewerState = "joined"
	// ViewerStateExited is mean that member has been exited broadcast
	ViewerStateExited ViewerState = "exited"
)

// Viewer struct
type Viewer struct {
	ID          int64       `json:"id" db:"id"`
	BroadcastID string      `json:"broadcast_id" db:"broadcast_id"`
	UserID      int64       `json:"user_id" db:"user_id"`
	State       ViewerState `json:"state" db:"state"`
	JoinedAt    time.Time   `json:"joined_at" db:"joined_at"`
	ExitedAt    time.Time   `json:"exited_at" db:"exited_at"`

	LocalSessionDescription  string `json:"local_sdp" db:"-"`
	RemoteSessionDescription string `json:"remote_sdp" db:"-"`

	SDPChan chan string `json:"-" db:"-"`

	db *sqlx.DB

	webrtc *Webrtc
}

// NewViewer creates new Viewer object with db conn, userID and broadcastID
func NewViewer(db *sqlx.DB, webrtc *Webrtc, userID int64, broadcastID string) *Viewer {
	return &Viewer{
		db:          db,
		UserID:      userID,
		BroadcastID: broadcastID,
		State:       ViewerStateJoined,
		webrtc:      webrtc,
		SDPChan:     make(chan string),
	}
}

// Create saves to db
func (v *Viewer) Create() error {
	insertQuery := `INSERT INTO viewers (broadcast_id, user_id, state, joined_at)
		VALUES ($1,$2,$3,NOW()) RETURNING id`
	return v.db.Get(&v.ID, insertQuery, v.BroadcastID, v.UserID, v.State)
}

// Run creates and runs main loop of viewer of broadcast
func (v *Viewer) Run(track *webrtc.Track) {
	go func(viewer *Viewer, track *webrtc.Track) {
		peerConnection, err := v.webrtc.NewPeerConnection()
		if err != nil {
			log.Println(err)
			return
		}

		if _, err = peerConnection.AddTrack(track); err != nil {
			log.Println(err)
			return
		}

		offer := webrtc.SessionDescription{}
		utils.DecodeSDP(viewer.LocalSessionDescription, &offer)

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

		log.Println("send SDP")
		viewer.SDPChan <- utils.EncodeSDP(answer)

		log.Println("Loop of viewer...")

		select {}
	}(v, track)
}
