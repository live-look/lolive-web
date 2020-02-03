package camforchat

import (
	"gitlab.com/isqad/camforchat/utils"

	"github.com/volatiletech/authboss"
	"go.uber.org/zap"
	"net/http"
)

// BroadcastsNew handles openning page for preparing to broadcast (preview of webcam and tools for start and stop
// broadcast)
func BroadcastsNew(w http.ResponseWriter, r *http.Request) {
	var data authboss.HTMLData
	dataIntf := r.Context().Value(authboss.CTXKeyData)

	if dataIntf == nil {
		data = authboss.HTMLData{}
	} else {
		data = dataIntf.(authboss.HTMLData)
	}

	logger, _ := GetLogger(r.Context())

	t, err := utils.Tmpl("layout.html", "broadcast.html")
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
