package event

import "github.com/netbrain/todoapp-go-es/common"

// Subscription handles a single subscription to a set of events from the event bus
type Subscription struct {
	Name      string
	EventChan chan *common.EventMessage
	eventType map[string]bool
	destroyed bool
}

// ChangeSubscription changes the current subscription to a new set of events
func (s *Subscription) ChangeSubscription(eventTypes ...string) {
	newMap := make(map[string]bool)
	for _, eventType := range eventTypes {
		newMap[eventType] = true
	}
	s.eventType = newMap
}

// Destroy marks this subscription for removal and closes it's event channel
func (s *Subscription) Destroy() {
	if !s.destroyed {
		s.destroyed = true
		close(s.EventChan)
	}
}
