package camforchat

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetWebrtcAPI(t *testing.T) {
	wrtc := NewWebrtc()
	ctx := context.WithValue(context.Background(), ctxKeyWebrtcAPI, wrtc)
	w, ok := GetWebrtcAPI(ctx)

	assert.True(t, ok)
	assert.IsType(t, NewWebrtc(), w)
}
