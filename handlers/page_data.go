package handlers

import (
	"camforchat/internal"
	"github.com/volatiletech/authboss"
)

// PageData is struct for passing data to views
type PageData struct {
	UserData         authboss.HTMLData
	Broadcasts       []*internal.Broadcast
	CurrentBroadcast *internal.Broadcast
}
