package todo

import (
	"encoding/json"
	"log"

	"github.com/netbrain/todoapp-go-es/common"
	"github.com/netbrain/todoapp-go-es/event"
	"github.com/netbrain/todoapp-go-es/fsstore"
)

//Projection the todo projection which creates todo views
type Projection struct {
	subscription *event.Subscription
	datastore    fsstore.FSStore
}

//NewProjection creates a new Projection
func NewProjection(bus event.Bus) *Projection {
	datastore, err := fsstore.NewJSONFSStore("todo")
	if err != nil {
		panic(err)
	}
	p := &Projection{
		subscription: bus.Subscribe(
			"TodoProjection",
			eventTodoItemCreated,
			eventTodoItemRemoved,
			eventTodoItemUpdated,
		),
		datastore: datastore,
	}

	go p.start()

	return p
}

//HandleEvent handles events this projection subscribes to
func (p *Projection) HandleEvent(event *common.EventMessage) {
	switch event.Name {
	case eventTodoItemUpdated:
		fallthrough
	case eventTodoItemCreated:
		p.handleTodoItemCreatedEvent(event)
	case eventTodoItemRemoved:
		p.handleTodoItemRemovedEvent(event)
	}
}

func (p *Projection) handleTodoItemCreatedEvent(event *common.EventMessage) {
	todo := new(Todo)
	err := json.Unmarshal(*event.Data, todo)
	if err != nil {
		panic(err)
	}
	p.datastore.Set(todo.ID, todo)
	p.datastore.AddToCollection("all", todo.ID, todo)
}

func (p *Projection) handleTodoItemRemovedEvent(event *common.EventMessage) {
	var id string
	err := json.Unmarshal(*event.Data, &id)
	if err != nil {
		log.Panic(err)
	}
	p.datastore.Remove(id)
	p.datastore.RemoveFromCollection("all", id)
}

func (p *Projection) start() {
	for {
		select {
		case event := <-p.subscription.EventChan:
			go p.HandleEvent(event)
		}
	}
}
