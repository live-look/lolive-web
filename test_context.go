package camforchat

import (
	"context"
	"github.com/jmoiron/sqlx"
)

type testContext struct {
	context context.Context
	db      *sqlx.DB
	webrtc  *Webrtc
}

func newTestContext(parent context.Context, db *sqlx.DB, webrtc *Webrtc) *testContext {
	return &testContext{context: parent, db: db, webrtc: webrtc}
}
