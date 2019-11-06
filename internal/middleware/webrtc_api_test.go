package middleware

import (
	"camforchat/internal/models"
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetWebrtcAPI(t *testing.T) {
	wrtc := models.NewWebrtc()
	ctx := context.WithValue(context.Background(), ctxKeyWebrtcAPI, wrtc)
	w, ok := GetWebrtcAPI(ctx)

	assert.True(t, ok)
	assert.IsType(t, models.NewWebrtc(), w)
}
