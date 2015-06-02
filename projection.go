package main

import "github.com/netbrain/todoapp-go-es/common"

type Projection interface {
	HandleEvent(*common.EventMessage)
}
