package models

type broadcastStorer interface {
	Save(*Broadcast) error
	Find(ID string) (*Broadcast, error)
	FindByState(s BroadcastState) ([]*Broadcast, error)
	UpdateState(ID string, s BroadcastState) error
}
