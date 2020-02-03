package camforchat

import (
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
)

// BroadcastsCreate handles creating broadcast
// TODO: extract into API
func BroadcastsCreate(w http.ResponseWriter, r *http.Request) {
	db, _ := GetDb(r.Context())
	logger, _ := GetLogger(r.Context())
	user, _ := GetCurrentUser(r.Context())
	webrtc, _ := GetWebrtcAPI(r.Context())

	broadcast, err := NewBroadcast(user.ID, user.Name, webrtc)
	if err != nil {
		logger.Error("initialize broadcast failed", zap.Error(err))
		http.Error(w, http.StatusText(500), 500)
		return
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(broadcast); err != nil {
		logger.Error("decoding request body failed", zap.Error(err))
		http.Error(w, http.StatusText(500), 500)
		return
	}

	if err := NewBroadcastDbStorage(db).Save(broadcast); err != nil {
		logger.Error("creating broadcast failed", zap.Error(err))
		http.Error(w, http.StatusText(422), 422)
		return
	}

	broadcastHandler, _ := GetBroadcastHandler(r.Context())
	broadcastHandler.StartBroadcasting(broadcast)

	remoteSdp := <-broadcast.SDPChan
	broadcast.RemoteSessionDescription = remoteSdp

	resp, err := json.Marshal(broadcast)
	if err != nil {
		logger.Error("marshaling broadacst failed", zap.Error(err))
		http.Error(w, http.StatusText(500), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}
