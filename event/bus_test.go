package event

import (
	"bytes"
	"testing"
	"time"

	"github.com/netbrain/todoapp-go-es/common"
)

func TestCanSubscribeToEvents(t *testing.T) {
	t.Parallel()
	var byteArray []byte
	buffer := bytes.NewBuffer(byteArray)
	eventBus := NewDefaultBus()
	repo := NewDefaultRepository(buffer, buffer, eventBus)
	subscriber := eventBus.Subscribe("test")

	go eventBus.Start()

	event := &common.EventMessage{
		Name: "test",
	}
	repo.Write(event)

	recvEvent := <-subscriber.EventChan
	if recvEvent != event {
		t.Fail()
	}
}

func TestCanQuitASubscribtion(t *testing.T) {
	t.Parallel()
	eventBus := NewDefaultBus()
	go eventBus.Start()

	subscriber := eventBus.Subscribe("test")
	subscriber.Destroy()

	if !subscriber.destroyed {
		t.Fatal("Should be destroyed")
	}

	select {
	case _, more := <-subscriber.EventChan:
		if more {
			t.Fatal("EventChan should be closed")
		}
	}

	subscriber = eventBus.Subscribe("test")

	var byteArray []byte
	buffer := bytes.NewBuffer(byteArray)
	repo := NewDefaultRepository(buffer, buffer, eventBus)
	event := &common.EventMessage{
		Name: "test",
	}
	repo.Write(event)

	recvEvent, more := <-subscriber.EventChan
	if recvEvent != event {
		t.Fatal(more)
	}
	time.Sleep(100)
	if len(eventBus.subscriptions) != 1 {
		t.Fatalf("Slice of subscribers should be 1, instead found %d", len(eventBus.subscriptions))
	}
}
