package camforchat

import (
	"github.com/volatiletech/authboss"
)

// PageData is struct for passing data to views
type PageData struct {
	UserData         authboss.HTMLData
	Broadcasts       []*Broadcast
	CurrentBroadcast *Broadcast
}
