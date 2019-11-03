package models

import (
	"context"
	"github.com/jmoiron/sqlx"
)

type testContext struct {
	context context.Context
	db      *sqlx.DB
}

func newTestContext(parent context.Context, db *sqlx.DB) *testContext {
	return &testContext{context: parent, db: db}
}
