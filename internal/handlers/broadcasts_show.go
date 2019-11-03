package handlers

import (
	appMiddleware "camforchat/internal/middleware"
	"camforchat/internal/usecases"
	"github.com/volatiletech/authboss"
	"go.uber.org/zap"
	"net/http"
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

	t, err := usecases.Tmpl("layout.html", "view.html")
	if err != nil {
		logger.Error("error during parse template", zap.Error(err))
		http.Error(w, http.StatusText(404), 404)
		return
	}

	err = t.Execute(w, &PageData{UserData: data})
	if err != nil {
		logger.Error("error during execute template", zap.Error(err))
		http.Error(w, http.StatusText(500), 500)
		return
	}
}
