package event

import (
	"bufio"
	"encoding/json"
	"io"
	"log"

	"github.com/netbrain/todoapp-go-es/common"
)

//Repository interface handles the writes and reads of the event log
type Repository interface {
	Write(*common.EventMessage) error
	Read() error
}

//DefaultRepository is the default implementation of a event repository
type DefaultRepository struct {
	r        io.Reader
	w        io.Writer
	eventBus Bus
}

//NewDefaultRepository instantiates a new DefaultRepository
func NewDefaultRepository(r io.Reader, w io.Writer, bus Bus) *DefaultRepository {
	return &DefaultRepository{
		r:        r,
		w:        w,
		eventBus: bus,
	}
}

func (d *DefaultRepository) Write(event *common.EventMessage) error {
	jsonEvent, err := json.Marshal(event)
	if err != nil {
		return err
	}

	if _, err := d.w.Write(append(jsonEvent, '\n')); err != nil {
		return err
	}

	json.Unmarshal(jsonEvent, event)
	d.eventBus.Notify(event)
	return nil
}

func (d *DefaultRepository) Read() error {
	scanner := bufio.NewScanner(d.r)
	for scanner.Scan() {
		event := &common.EventMessage{}
		if err := json.Unmarshal(scanner.Bytes(), event); err != nil {
			return err
		}
		log.Printf("Event: %s, version: %d", event.Name, event.Version)
		d.eventBus.Notify(event)
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}
