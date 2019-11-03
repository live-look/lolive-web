package handlers

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/volatiletech/authboss"
	"go.uber.org/zap"
	"net/http"
	"strconv"

	appMiddleware "camforchat/internal/middleware"
	"camforchat/internal/models"
)

// ViewersCreate handles POST /broadcasts/{broadcastID}/viewers
// creates new Viewer and runs it
func ViewersCreate(w http.ResponseWriter, r *http.Request) {
	logger, _ := appMiddleware.GetLog(r.Context())

	broadcastID, err := strconv.ParseInt(chi.URLParam(r, "broadcastID"), 10, 64)
	if err != nil {
		logger.Error("parse broadcastID failed", zap.Error(err))
		http.Error(w, http.StatusText(400), 400)
		return
	}

	dataIntf := r.Context().Value(authboss.CTXKeyData)
	if dataIntf == nil {
		http.Error(w, http.StatusText(401), 401)
		return
	}

	data := dataIntf.(authboss.HTMLData)

	db, _ := appMiddleware.GetDb(r.Context())
	broadcast, err := models.FindBroadcast(db, broadcastID)
	if err != nil {
		http.Error(w, http.StatusText(404), 404)
		return
	}

	viewer := models.NewViewer(db, data["current_user_id"].(int64), broadcast.ID)
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
