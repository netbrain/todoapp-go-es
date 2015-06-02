package event

import (
	"log"
	"time"

	"github.com/netbrain/todoapp-go-es/common"
)

//Bus interface
type Bus interface {
	Notify(*common.EventMessage)
	Subscribe(string, ...string) *Subscription
	Start()
}

//DefaultBus implementation
type DefaultBus struct {
	subscriptions []*Subscription
	notifyChan    chan *common.EventMessage
}

// NewDefaultBus creates a new default event bus
func NewDefaultBus() *DefaultBus {
	return &DefaultBus{
		subscriptions: make([]*Subscription, 0),
		notifyChan:    make(chan *common.EventMessage, 0),
	}

}

//Subscribe listeners to the event bus
func (d *DefaultBus) Subscribe(name string, eventType ...string) *Subscription {
	eventTypeMap := make(map[string]bool)
	for _, v := range eventType {
		eventTypeMap[v] = true
	}

	subscription := &Subscription{
		Name:      name,
		EventChan: make(chan *common.EventMessage, 1),
		eventType: eventTypeMap,
	}
	d.subscriptions = append(d.subscriptions, subscription)
	return subscription
}

//Notify listeners of a new event
func (d *DefaultBus) Notify(event *common.EventMessage) {
	d.notifyChan <- event
}

//Start should be run in it's own go routine, this starts the eventbus for relaying events to listeners
func (d *DefaultBus) Start() {
	for {
		select {
		case event := <-d.notifyChan:
			for i := len(d.subscriptions) - 1; i >= 0; i-- {
				li := i
				subscription := d.subscriptions[li]
				if subscription.destroyed {
					d.subscriptions = append(d.subscriptions[:li], d.subscriptions[li+1:]...)
				} else if len(subscription.eventType) == 0 || subscription.eventType[event.Name] {
					go func() {
						select {
						case subscription.EventChan <- event:
							log.Printf("Sending event %s to %s", event.Name, subscription.Name)
						case <-time.After(3 * time.Second):
							log.Printf("Sending event to %s timed out!", subscription.Name)
						}
					}()
				}
			}
		}
	}
}
