package main

import (
	"fmt"

	"github.com/netbrain/todoapp-go-es/common"

	"testing"
)

const (
	cmdName = "someCommand"
)

func TestCanRegisterAndHandleACommand(t *testing.T) {
	t.Parallel()
	cmdHandler := NewDefaultCommandHandler()
	hasRun := false
	var handler CommandFunc = func(*common.CommandMessage, chan<- *common.EventMessage) error {
		hasRun = true
		return nil
	}
	if err := cmdHandler.RegisterCommand(cmdName, handler); err != nil {
		t.Fatal(err)
	}

	if err := cmdHandler.HandleCommandMessage(&common.CommandMessage{
		Name: cmdName,
	}); err != nil {
		t.Fatal(err)
	}

	if !hasRun {
		t.Fail()
	}
}

func TestCannotRegisterACommandThatAlreadyExists(t *testing.T) {
	t.Parallel()
	cmdHandler := NewDefaultCommandHandler()
	if err := cmdHandler.RegisterCommand(cmdName, nil); err != nil {
		t.Fatal(err)
	}

	if err := cmdHandler.RegisterCommand(cmdName, nil); err == nil {
		t.Fatal("No error when trying to register a commnand that already exists!")
	}
}

func TestHandlersFunctionsAsMiddleware(t *testing.T) {
	t.Parallel()
	cmdHandler := NewDefaultCommandHandler()
	middlewaresRun := 0

	var handlers = []CommandFunc{
		func(*common.CommandMessage, chan<- *common.EventMessage) error {

			middlewaresRun++
			return nil
		},
		func(*common.CommandMessage, chan<- *common.EventMessage) error {
			middlewaresRun++
			return fmt.Errorf("Some error occured at middleware no 2")

		},
		func(*common.CommandMessage, chan<- *common.EventMessage) error {
			middlewaresRun++
			return nil
		},
	}

	if err := cmdHandler.RegisterCommand(cmdName, handlers...); err != nil {
		t.Fatal(err)
	}

	if err := cmdHandler.HandleCommandMessage(&common.CommandMessage{
		Name: cmdName,
	}); err == nil {
		t.Fatal(err)
	}

	if middlewaresRun != 2 {
		t.Fail()
	}

}
