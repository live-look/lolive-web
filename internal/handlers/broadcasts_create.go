package handlers

import (
	"encoding/json"
	"go.uber.org/zap"
	"net/http"

	appMiddleware "camforchat/internal/middleware"
	"camforchat/internal/models"
)

// BroadcastsCreate handles creating broadcast
func BroadcastsCreate(w http.ResponseWriter, r *http.Request) {
	db, _ := appMiddleware.GetDb(r.Context())
	logger, _ := appMiddleware.GetLog(r.Context())
	user, _ := appMiddleware.GetCurrentUser(r.Context())

	broadcast, err := models.CreateBroadcast(db, user)
	if err != nil {
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
