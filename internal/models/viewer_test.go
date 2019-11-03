package models

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestViewerCreate(t *testing.T) {
	ctx := setupTestContext(t)
	defer ctx.db.Close()

	user := userFixture(t, ctx, "viewer_creating@camforchat-test.net")
	broadcast := broadcastFixture(t, ctx, user)

	v := NewViewer(ctx.db, user.ID, broadcast.ID)
	err := v.Create()
	assert.Nil(t, err)
}
