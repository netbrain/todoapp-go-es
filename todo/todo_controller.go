package todo

import (
	"encoding/json"

	"code.google.com/p/go-uuid/uuid"

	"github.com/netbrain/todoapp-go-es/common"
)

const (
	eventTodoItemCreated = "todoItemCreated"
	eventTodoItemRemoved = "todoItemRemoved"
	eventTodoItemUpdated = "todoItemUpdated"
)

//CreateTodoItem creates a todo based on a command message
func CreateTodoItem(cmd *common.CommandMessage, eventChan chan<- *common.EventMessage) error {
	var todo Todo
	if err := json.Unmarshal(*cmd.Data, &todo); err != nil {
		return err
	}
	todo.ID = uuid.New()

	data, err := json.Marshal(todo)
	if err != nil {
		return err
	}

	raw := json.RawMessage(data)

	event := &common.EventMessage{
		Name: eventTodoItemCreated,
		Data: &raw,
	}
	eventChan <- event
	return nil
}

//RemoveTodoItem removes a todo based on a command message
func RemoveTodoItem(cmd *common.CommandMessage, eventChan chan<- *common.EventMessage) error {
	event := &common.EventMessage{
		Name: eventTodoItemRemoved,
		Data: cmd.Data,
	}
	eventChan <- event
	return nil
}

//UpdateTodoItem updates a todo based on a command message
func UpdateTodoItem(cmd *common.CommandMessage, eventChan chan<- *common.EventMessage) error {
	event := &common.EventMessage{
		Name: eventTodoItemUpdated,
		Data: cmd.Data,
	}
	eventChan <- event
	return nil
}
