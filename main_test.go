package main

import (
	"net/http/httptest"
	"sync"
)

var serverAddr string
var once sync.Once

func startServer() {
	//initServer()
	server := httptest.NewServer(nil)
	serverAddr = server.Listener.Addr().String()
}

/*func TestCreateTodo(t *testing.T) {
	once.Do(startServer)

	origin := fmt.Sprintf("http://%s", serverAddr)
	url := fmt.Sprintf("ws://%s/cmd", serverAddr)

	go func() {
		cmdWs, err := websocket.Dial(url, "", origin)
		if err != nil {
			t.Fatal(err)
		}

		if err := websocket.JSON.Send(cmdWs, &common.CommandMessage{
			Name:   "someCommand",
			Params: map[string]string{"Name": "My TODO"},
		}); err != nil {
			t.Fatal(err)
		}

		cmdWs.Close()
	}()

	url = fmt.Sprintf("ws://%s/stream", serverAddr)
	streamWs, err := websocket.Dial(url, "", origin)
	if err != nil {
		t.Fatal(err)
	}

	var msg = make([]byte, 512)
	var n int
	if n, err = streamWs.Read(msg); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Received: %s.\n", msg[:n])
	streamWs.Close()

}*/
