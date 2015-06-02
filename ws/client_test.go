package ws

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sync"
	"testing"

	"github.com/netbrain/todoapp-go-es/common"

	"github.com/netbrain/todoapp-go-es/event"

	"golang.org/x/net/websocket"
)

var clientConn *websocket.Conn
var wsClient *Client
var once sync.Once

var eventBus = event.NewDefaultBus()

func init() {
	go eventBus.Start()
	log.SetFlags(log.Flags() | log.Lshortfile)
}

func startServerAndClient() {
	var err error
	http.Handle("/ws", websocket.Handler(wsHandler))
	server := httptest.NewServer(nil)
	serverAddr := server.Listener.Addr().String()

	origin := "http://localhost/"
	url := fmt.Sprintf("ws://%s/ws", serverAddr)
	clientConn, err = websocket.Dial(url, "", origin)
	if err != nil {
		panic(err)
	}
}

func wsHandler(c *websocket.Conn) {
	wsClient = NewClient(c, eventBus)
	wsClient.Listen()
}

func TestClient(t *testing.T) {
	once.Do(startServerAndClient)

	sendEvent := &common.EventMessage{Name: "test"}
	eventBus.Notify(sendEvent)

	recvEvent := &common.EventMessage{}
	//ignore first
	websocket.JSON.Receive(clientConn, nil)
	websocket.JSON.Receive(clientConn, recvEvent)

	if !reflect.DeepEqual(sendEvent, recvEvent) {
		t.Fatalf("Events isnt equal, %#v != %#v", sendEvent, recvEvent)
	}
}
