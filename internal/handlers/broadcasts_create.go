package handlers

import (
	"encoding/json"
	"github.com/volatiletech/authboss"
	"go.uber.org/zap"
	"net/http"

	appMiddleware "camforchat/internal/middleware"
	"camforchat/internal/models"
)

// BroadcastsCreate handles creating broadcast
func BroadcastsCreate(w http.ResponseWriter, r *http.Request) {
	dataIntf := r.Context().Value(authboss.CTXKeyData)
	if dataIntf == nil {
		http.Error(w, http.StatusText(401), 401)
		return
	}

	data := dataIntf.(authboss.HTMLData)

	db, _ := appMiddleware.GetDb(r.Context())
	logger, _ := appMiddleware.GetLog(r.Context())

	broadcast := models.NewBroadcast(db, data["current_user_id"].(int64))
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(broadcast); err != nil {
		logger.Error("decoding request body failed", zap.Error(err))
		http.Error(w, http.StatusText(500), 500)
		return
	}

	if err := broadcast.Create(); err != nil {
		logger.Error("creating broadcast failed", zap.Error(err))
		http.Error(w, http.StatusText(422), 422)
		return
	}

	broadcastHandler, _ := appMiddleware.GetBroadcastHandler(r.Context())
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