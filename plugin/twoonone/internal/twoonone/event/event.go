package event

import (
	"context"
	"sync"
	"time"

	twoonone_pb "github.com/nanachi-sh/susubot-code/plugin/twoonone/pkg/protos/twoonone"
)

type EventStream interface {
	Read() (*twoonone_pb.RoomEventResponse, bool) //error is *types.AppError
	Emit(*twoonone_pb.RoomEventResponse)
	Close()
}

type eventStream struct {
	events []*twoonone_pb.RoomEventResponse

	roomHash string
	read     sync.RWMutex
	wait     sync.RWMutex
	now      context.Context
	close    bool
}

type myevent struct{}

func NewEventStream(roomHash string) EventStream {
	es := &eventStream{
		roomHash: roomHash,
		read:     sync.RWMutex{},
		wait:     sync.RWMutex{},
		now:      context.Background(),
	}
	es.read.Lock()
	go es.emits()
	return es
}

func (e *eventStream) Read() (*twoonone_pb.RoomEventResponse, bool) {
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
	if e.close {
		return nil, false
	}
	if v, ok := e.now.Value(myevent{}).(*twoonone_pb.RoomEventResponse); !ok {
		return nil, false
	} else {
		return v, true
	}
}

func (e *eventStream) Emit(resp *twoonone_pb.RoomEventResponse) {
	if e.close {
		return
	}
	// 空队列，且上一个emit已处理完
	if e.now.Value(myevent{}) == nil && len(e.events) == 0 {
		e.emit(resp)
	} else {
		e.events = append(e.events, resp)
	}
}

func (e *eventStream) emit(resp *twoonone_pb.RoomEventResponse) {
	e.now = context.WithValue(e.now, myevent{}, resp)
	e.read.Unlock()
}

func (e *eventStream) emits() {
	for {
		time.Sleep(time.Millisecond)
		if len(e.events) == 0 {
			time.Sleep(time.Millisecond * 20)
			continue
		}
		if e.now != nil {
			continue
		}
		resp := e.events[len(e.events)-1]
		if len(e.events) == 1 {
			e.events = []*twoonone_pb.RoomEventResponse{}
		} else {
			e.events = e.events[:len(e.events)-1]
		}
		e.emit(resp)
	}
}

func (e *eventStream) Close() {
	if e.close {
		return
	}
	e.close = true
	e.read.TryLock()
	e.read.Unlock()
	e.wait.TryLock()
	e.wait.Unlock()
}

func (e *eventStream) clean() {
	e.now = context.WithValue(e.now, myevent{}, nil)
}
