package connect

import (
	"errors"
	"sync"
)

type Container struct {
	Buckets      []*Bucket
	BucketNumber int
}

func NewContainer(lenBucket int) *Container {
	c := new(Container)
	c.Buckets = make([]*Bucket, lenBucket)
	c.BucketNumber = lenBucket
	for i := 0; i < lenBucket; i++ {
		c.Buckets[i] = NewBucket()
	}
	return c
}

type Bucket struct {
	chMap   map[int]*Channel
	roomMap map[int]*Room
	Lock    sync.RWMutex
}

func NewBucket() *Bucket {
	b := new(Bucket)
	b.roomMap = make(map[int]*Room, 10)
	b.chMap = make(map[int]*Channel, 100)
	return b
}

func (b *Bucket) Channel(uid int) (*Channel, error) {
	b.Lock.RLock()
	c, ok := b.chMap[uid]
	b.Lock.RUnlock()
	if !ok {
		return nil, errors.New("the user not have channel")
	}
	return c, nil
}

func (b *Bucket) Room(RoomId int) (*Room, error) {
	b.Lock.RLock()
	r, ok := b.roomMap[RoomId]
	b.Lock.RUnlock()
	if !ok {
		return nil, errors.New("the room is not exit")
	}
	return r, nil
}

func (b *Bucket) Push(uid int, roomId int, ch *Channel) error {
	var (
		room *Room
		ok   bool
	)
	b.Lock.Lock()
	if roomId != 0 {
		room, ok = b.roomMap[roomId]
		if !ok {
			room = new(Room)
			b.roomMap[roomId] = room
		}
	}
	b.chMap[uid] = ch
	b.Lock.Unlock()
	if room != nil {
		err := room.Put(ch)
		return err
	}
	return nil
}

func (b *Bucket) DeleteChannel(uid int) {
	b.Lock.Lock()
	delete(b.chMap, uid)
	b.Lock.Unlock()
}

func (b *Bucket) DeleteRoom(roomId int) {
	b.Lock.Lock()
	delete(b.roomMap, roomId)
	b.Lock.Unlock()
}
