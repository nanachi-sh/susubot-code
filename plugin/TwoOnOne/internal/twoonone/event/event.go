package event

import (
	"context"
	"errors"
	"sync"

	twoonone_pb "github.com/nanachi-sh/susubot-code/plugin/twoonone/pkg/protos/twoonone"
)

var (
	events []*EventStream
)

type EventStream struct {
	roomHash string
	read     sync.RWMutex
	wait     sync.RWMutex
	now      context.Context
}

type myevent struct{}

func NewEventStream(roomHash string) *EventStream {
	es := &EventStream{
		roomHash: roomHash,
		read:     sync.RWMutex{},
		wait:     sync.RWMutex{},
		now:      context.Background(),
	}
	es.read.Lock()
	return es
}

func FindEventStream(roomHash string) (*EventStream, bool) {
	for _, v := range events {
		if v.roomHash == roomHash {
			return v, true
		}
	}
	return nil, false
}

func (e *EventStream) Read() (*twoonone_pb.EventRoomResponse, error) {
	e.wait.RLock()
	e.wait.RUnlock()
	e.read.RLock()
	defer func() {
		if e.read.TryLock() {
			// 处理
			e.clean()
			e.wait.Unlock()
		}
	}()
	defer e.read.RUnlock()
	e.wait.TryLock()
	if v, ok := e.now.Value(myevent{}).(*twoonone_pb.EventRoomResponse); !ok {
		return nil, errors.New("异常错误")
	} else {
		return v, nil
	}
}

func (e *EventStream) Emit(resp *twoonone_pb.EventRoomResponse) {
	e.wait.RLock()
	e.wait.RUnlock()
	e.now = context.WithValue(e.now, myevent{}, resp)
	e.read.Unlock()
}

func (e *EventStream) clean() {
	e.now = context.WithValue(e.now, myevent{}, nil)
}
