package handlers

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
	"net/http"
	"strconv"

	"camforchat/internal"
	appMiddleware "camforchat/internal/middleware"
)

// ViewersCreate handles POST /broadcasts/{broadcastID}/viewers
// creates new Viewer and runs it
// TODO: extract API
func ViewersCreate(w http.ResponseWriter, r *http.Request) {
	logger, _ := appMiddleware.GetLog(r.Context())

	broadcastID, err := strconv.ParseInt(chi.URLParam(r, "broadcastID"), 10, 64)
	if err != nil {
		logger.Error("parse broadcastID failed", zap.Error(err))
		http.Error(w, http.StatusText(400), 400)
		return
	}

	db, _ := appMiddleware.GetDb(r.Context())
	webrtc, _ := appMiddleware.GetWebrtcAPI(r.Context())
	broadcast, err := internal.FindBroadcast(db, webrtc, broadcastID)
	if err != nil {
		http.Error(w, http.StatusText(404), 404)
		return
	}

	user, _ := appMiddleware.GetCurrentUser(r.Context())

	viewer := internal.NewViewer(db, webrtc, user.ID, broadcast.ID)
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

	broadcastHandler, _ := appMiddleware.GetBroadcastHandler(r.Context())
	logger.Info("start viewer")
	broadcastHandler.StartView(viewer)

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
