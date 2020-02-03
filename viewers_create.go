package camforchat

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
	"net/http"
)

// ViewersCreate handles POST /broadcasts/{broadcastID}/viewers
// creates new Viewer and runs it
// TODO: extract API
func ViewersCreate(w http.ResponseWriter, r *http.Request) {
	logger, _ := GetLogger(r.Context())

	broadcastID := chi.URLParam(r, "broadcastID")

	db, _ := GetDb(r.Context())
	webrtc, _ := GetWebrtcAPI(r.Context())
	broadcast, err := NewBroadcastDbStorage(db).Find(broadcastID)
	if err != nil {
		http.Error(w, http.StatusText(404), 404)
		return
	}

	user, _ := GetCurrentUser(r.Context())

	viewer := NewViewer(db, webrtc, user.ID, broadcast.ID)
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(viewer); err != nil {
		logger.Error("decoding request body failed", zap.Error(err))
		http.Error(w, http.StatusText(500), 500)
		return
	}

	if err := viewer.Create(); err != nil {
		logger.Error("creating viewer failed", zap.Error(err))
		http.Error(w, http.StatusText(422), 422)
		return
	}

	//broadcastHandler, _ := GetBroadcastHandler(r.Context())
	logger.Info("start viewer")
	//broadcastHandler.StartView(viewer)

	remoteSdp := <-viewer.SDPChan
	viewer.RemoteSessionDescription = remoteSdp

	resp, err := json.Marshal(viewer)
	if err != nil {
		logger.Error("marshaling viewer failed", zap.Error(err))
		http.Error(w, http.StatusText(500), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}
