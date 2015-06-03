package main

import (
	"net/http"
	"strings"

	"github.com/netbrain/todoapp-go-es/fsstore"
)

type QueryHandler struct {
	datastore fsstore.FSStore
}

func NewQueryHandler(datastore fsstore.FSStore) *QueryHandler {
	return &QueryHandler{
		datastore: datastore,
	}
}

func (q *QueryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		reqArr := strings.Split(r.RequestURI, "/")
		id := reqArr[len(reqArr)-1]

		data, err := todoDataStore.GetBytes(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		data.WriteTo(w)
	} else {
		http.NotFound(w, r)
		return
	}
}
