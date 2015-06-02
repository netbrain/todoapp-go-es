package ws

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"sync/atomic"

	"github.com/netbrain/todoapp-go-es/common"
	"github.com/netbrain/todoapp-go-es/event"

	"golang.org/x/net/websocket"
)

type Client struct {
	conn         *websocket.Conn
	eventBus     event.Bus
	subscription *event.Subscription
	quitChan     chan bool
}

const eventNewWebSocketClient = "newWebSocketClient"

var numClient = int32(0)

type subscribe struct {
	events []string `json:"events"`
}

func NewClient(conn *websocket.Conn, eventBus event.Bus) *Client {

	c := &Client{
		conn:         conn,
		eventBus:     eventBus,
		subscription: eventBus.Subscribe(fmt.Sprint("WS: %s", conn.LocalAddr())),
		quitChan:     make(chan bool, 1),
	}

	return c
}

func (c *Client) sendNumClientsEvent() {
	jsonData, err := json.Marshal(numClient)
	if err != nil {
		panic(err)
	}

	rawData := json.RawMessage(jsonData)

	c.eventBus.Notify(&common.EventMessage{
		Name: eventNewWebSocketClient,
		Data: &rawData,
	})
}

func (c *Client) Listen() {
	atomic.AddInt32(&numClient, 1)
	c.sendNumClientsEvent()

	go c.listenWrite()
	c.listenRead()

	c.subscription.Destroy()
	atomic.AddInt32(&numClient, -1)
	c.sendNumClientsEvent()
}

func (c *Client) listenRead() {
	s := new(subscribe)
	for {
		select {
		case <-c.quitChan:
			c.Stop()
			return
		default:
			err := websocket.JSON.Receive(c.conn, s)
			if err == io.EOF {
				c.Stop()
			} else if err != nil {
				c.Stop()
				log.Fatal(err)
			} else {
				log.Printf("Changing subscription to: %v", s.events)
				c.subscription.ChangeSubscription(s.events...)
			}
		}
	}
}

func (c *Client) listenWrite() {
	for {
		select {
		case <-c.quitChan:
			c.Stop()
			return
		case message := <-c.subscription.EventChan:
			err := websocket.JSON.Send(c.conn, message)
			if err != nil {
				c.Stop()
				log.Fatal(err)
			}
		}
	}
}

func (c *Client) Stop() {
	c.quitChan <- true
}
