package fsstore

import (
	"bytes"
	"os"
)

//FSStore ...
type FSStore interface {
	Set(string, interface{}) error
	Remove(string) error
	Get(string, interface{}) error
	GetBytes(id string) (*bytes.Reader, error)
	AddToCollection(string, string, interface{}) error
	RemoveFromCollection(string, string) error
	GetCollection(string) (map[string]interface{}, error)
	RemoveAll() error
	GetDataDir() string
}

//DataDir sets the root dir for data storage for this package
var DataDir = os.TempDir()
