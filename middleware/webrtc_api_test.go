package middleware

import (
	"camforchat/internal"
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetWebrtcAPI(t *testing.T) {
	wrtc := internal.NewWebrtc()
	ctx := context.WithValue(context.Background(), ctxKeyWebrtcAPI, wrtc)
	w, ok := GetWebrtcAPI(ctx)

	assert.True(t, ok)
	assert.IsType(t, internal.NewWebrtc(), w)
}
