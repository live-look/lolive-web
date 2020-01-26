package models

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestSaveBroadcast(t *testing.T) {
	ctx := setupTestContext(t)
	defer ctx.db.Close()

	user := userFixture(t, ctx, "broadcast_creating@camforchat-test.net")

	broadcast := broadcastFixture(user)

	bs := NewBroadcastDbStorage(ctx.db)

	err := bs.Save(broadcast)
	assert.Nil(t, err)
}

func TestFindBroadcast(t *testing.T) {
	ctx := setupTestContext(t)
	defer ctx.db.Close()

	user := userFixture(t, ctx, "broadcast_finding@camforchat-test.net")
	broadcast := broadcastFixture(user)

	bs := NewBroadcastDbStorage(ctx.db)

	err := bs.Save(broadcast)
	assert.Nil(t, err)

	b, err := bs.Find(broadcast.ID)
	assert.Nil(t, err)

	assert.Equal(t, user.ID, b.UserID)
	assert.Equal(t, user.Name, b.UserName)
	assert.Equal(t, broadcast.State, b.State)
	assert.Equal(t, broadcast.ID, b.ID)

	assert.NotNil(t, broadcast.Publish)
	assert.NotNil(t, broadcast.SDPChan)
}
