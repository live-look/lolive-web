package handlers

import (
	appMiddleware "camforchat/internal/middleware"
	"camforchat/internal/models"
	"camforchat/internal/usecases"
	"github.com/go-chi/chi"
	"github.com/volatiletech/authboss"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

// BroadcastsShow handles view of runned broadcast
func BroadcastsShow(w http.ResponseWriter, r *http.Request) {
	var data authboss.HTMLData
	dataIntf := r.Context().Value(authboss.CTXKeyData)

	if dataIntf == nil {
		data = authboss.HTMLData{}
	} else {
		data = dataIntf.(authboss.HTMLData)
	}

	logger, _ := appMiddleware.GetLog(r.Context())

	broadcastID, err := strconv.ParseInt(chi.URLParam(r, "broadcastID"), 10, 64)
	if err != nil {
		logger.Error("parse broadcastID failed", zap.Error(err))
		http.Error(w, http.StatusText(400), 400)
		return
	}

	db, _ := appMiddleware.GetDb(r.Context())
	broadcast, err := models.FindBroadcast(db, broadcastID)
	if err != nil {
		http.Error(w, http.StatusText(404), 404)
		return
	}

	t, err := usecases.Tmpl("layout.html", "view.html")
	if err != nil {
		logger.Error("error during parse template", zap.Error(err))
		http.Error(w, http.StatusText(404), 404)
		return
	}

	err = t.Execute(w, &PageData{UserData: data, CurrentBroadcast: broadcast})
	if err != nil {
		logger.Error("error during execute template", zap.Error(err))
		http.Error(w, http.StatusText(500), 500)
		return
	}
}
