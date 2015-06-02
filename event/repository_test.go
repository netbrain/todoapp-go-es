package event

import (
	"bytes"
	"reflect"
	"testing"
	"time"

	"github.com/netbrain/todoapp-go-es/common"
)

func TestCanWriteAndReadEvent(t *testing.T) {
	t.Parallel()
	buffer := new(bytes.Buffer)
	bus := NewDefaultBus()
	go bus.Start()
	repo := NewDefaultRepository(buffer, buffer, bus)
	event := &common.EventMessage{
		Name: "test",
	}
	repo.Write(event)

	subscripiton := bus.Subscribe("test")

	if err := repo.Read(); err != nil {
		t.Fatal(err)
	}

	select {
	case readEvent := <-subscripiton.EventChan:
		if !reflect.DeepEqual(readEvent, event) {
			t.Fatal("readEvent not equal to event")
		}
	case <-time.After(time.Millisecond * 100):
		t.Fatal("Timed out")
	}

}
