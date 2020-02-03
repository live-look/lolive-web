package camforchat

import (
	"github.com/go-chi/chi"
	"github.com/volatiletech/authboss"
	"gitlab.com/isqad/camforchat/utils"
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

	logger, _ := GetLogger(r.Context())

	broadcastID := chi.URLParam(r, "broadcastID")

	db, _ := GetDb(r.Context())
	broadcast, err := NewBroadcastDbStorage(db).Find(broadcastID)
	if err != nil {
		http.Error(w, http.StatusText(404), 404)
		return
	}

	t, err := utils.Tmpl("layout.html", "view.html")
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
