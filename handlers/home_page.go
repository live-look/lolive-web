package handlers

import (
	"camforchat/internal"
	appMiddleware "camforchat/internal/middleware"
	"camforchat/internal/usecases"
	"github.com/volatiletech/authboss"
	"go.uber.org/zap"
	"net/http"
)

// HomePage handles main page of site
func HomePage(w http.ResponseWriter, r *http.Request) {
	var data authboss.HTMLData
	dataIntf := r.Context().Value(authboss.CTXKeyData)

	if dataIntf == nil {
		data = authboss.HTMLData{}
	} else {
		data = dataIntf.(authboss.HTMLData)
	}

	logger, _ := appMiddleware.GetLog(r.Context())
	db, _ := appMiddleware.GetDb(r.Context())

	broadcasts, err := internal.GetBroadcastsByState(db, internal.BroadcastStateOnline)
	if err != nil {
		logger.Error("error during retreive broadcasts", zap.Error(err))
		http.Error(w, http.StatusText(500), 500)
		return
	}

	t, err := usecases.Tmpl("layout.html", "index.html")
	if err != nil {
		logger.Error("error during parse template", zap.Error(err))
		http.Error(w, http.StatusText(404), 404)
		return
	}

	err = t.Execute(w, &PageData{UserData: data, Broadcasts: broadcasts})
	if err != nil {
		logger.Error("error during execute template", zap.Error(err))
		http.Error(w, http.StatusText(500), 500)
		return
	}
}
