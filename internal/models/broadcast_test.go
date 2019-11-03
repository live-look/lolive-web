package models

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateBroadcast(t *testing.T) {
	ctx := setupTestContext(t)
	defer ctx.db.Close()

	user := userFixture(t, ctx, "broadcast_creating@camforchat-test.net")
	broadcast, err := CreateBroadcast(ctx.db, user)
	assert.Nil(t, err)

	assert.Equal(t, BroadcastStateOffline, broadcast.State)
	assert.Equal(t, user.ID, broadcast.UserID)
	assert.Equal(t, user.Name, broadcast.UserName)
	assert.NotEqual(t, 0, broadcast.ID)

	assert.NotNil(t, broadcast.Publish)
	assert.NotNil(t, broadcast.SDPChan)
}

func TestFindBroadcast(t *testing.T) {
	ctx := setupTestContext(t)
	defer ctx.db.Close()

	user := userFixture(t, ctx, "broadcast_finding@camforchat-test.net")
	broadcast := broadcastFixture(t, ctx, user)

	b, err := FindBroadcast(ctx.db, broadcast.ID)
	assert.Nil(t, err)

	assert.Equal(t, user.ID, b.UserID)
	assert.Equal(t, user.Name, b.UserName)
	assert.Equal(t, broadcast.State, b.State)
	assert.Equal(t, broadcast.ID, b.ID)

	assert.NotNil(t, broadcast.Publish)
	assert.NotNil(t, broadcast.SDPChan)
}
