package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"golang.org/x/net/websocket"

	"github.com/netbrain/todoapp-go-es/ws"

	"github.com/netbrain/todoapp-go-es/common"
	"github.com/netbrain/todoapp-go-es/event"
	"github.com/netbrain/todoapp-go-es/fsstore"
	"github.com/netbrain/todoapp-go-es/todo"
)

var eventBus event.Bus
var eventLogFile string
var eventLogWriter io.Writer
var eventLogReader io.Reader
var eventRepository event.Repository
var todoProjection *todo.Projection
var todoDataStore fsstore.FSStore
var staticPath string

func init() {
	//log.SetOutput(ioutil.Discard)
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	staticPath = filepath.Join(dir, "static")
	fsstore.DataDir = filepath.Join(staticPath, "api")
	eventBus = event.NewDefaultBus()
	eventLogFile = filepath.Join(os.TempDir(), "eventlog")

	todoDataStore, err = fsstore.NewJSONFSStore("todo")
	if err != nil {
		log.Fatal(err)
	}

	todoProjection = todo.NewProjection(eventBus, todoDataStore)
	eventLogWriter, _ = os.OpenFile(eventLogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	eventLogReader, _ = os.Open(eventLogFile)
	eventRepository = event.NewDefaultRepository(eventLogReader, eventLogWriter, eventBus)

	log.SetFlags(log.Flags() | log.Lshortfile)
}

func main() {
	cmdHandler := NewDefaultCommandHandler()
	cmdHandler.RegisterCommand("createTodoItem", todo.CreateTodoItem)
	cmdHandler.RegisterCommand("removeTodoItem", todo.RemoveTodoItem)
	cmdHandler.RegisterCommand("updateTodoItem", todo.UpdateTodoItem)
	go cmdHandler.Start()
	go eventBus.Start()

	//Read the event log
	go func() {
		log.Println("Reading the event log...")
		err := eventRepository.Read()
		if err != nil {
			panic(err)
		}
	}()

	http.HandleFunc("/cmd/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			data, err := ioutil.ReadAll(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
			}

			cmd := new(common.CommandMessage)
			err = json.Unmarshal(data, cmd)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
			}

			if err := cmdHandler.HandleCommandMessage(cmd); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	})
	http.Handle("/api/todo/", NewQueryHandler(todoDataStore))
	http.Handle("/", http.FileServer(http.Dir(filepath.Join(staticPath, "app"))))
	http.Handle("/ws/", websocket.Handler(wsHandler))

	log.Println("Listening...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}

func wsHandler(conn *websocket.Conn) {
	log.Println("New WS client")
	ws := ws.NewClient(conn, eventBus)
	ws.Listen()
}
