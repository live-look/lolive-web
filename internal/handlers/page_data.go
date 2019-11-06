package handlers

import (
	"camforchat/internal/models"
	"github.com/volatiletech/authboss"
)

// PageData is struct for passing data to views
type PageData struct {
	UserData         authboss.HTMLData
	Broadcasts       []*models.Broadcast
	CurrentBroadcast *models.Broadcast
}
