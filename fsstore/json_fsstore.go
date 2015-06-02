package fsstore

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"
)

//JSONFSStore a store that encodes data to json and gzips it
type JSONFSStore struct {
	dataDir string
	files   map[string]*jsonFile
}

//NewJSONFSStore creates a new JSONFSStore with a relative or absolute datadir path
func NewJSONFSStore(dataDir string) (*JSONFSStore, error) {
	if dataDir[0] != '/' {
		dataDir = filepath.Join(DataDir, dataDir)
	}

	store := &JSONFSStore{
		dataDir: dataDir,
		files:   make(map[string]*jsonFile),
	}

	store.RemoveAll()
	return store, nil
}

//Set sets the data assosciated with a given id
func (j *JSONFSStore) Set(id string, data interface{}) error {
	file, err := j.getJsonFile(id)
	if err != nil {
		return err

	}
	return file.set(data)
}

//Remove removes the data associated with a given id
func (j *JSONFSStore) Remove(id string) error {
	file, err := j.getJsonFile(id)
	if err != nil {
		return err
	}

	file.remove()
	return nil
}

//Get injects the data assosciated with a given id
func (j *JSONFSStore) Get(id string, v interface{}) error {
	file, err := j.getJsonFile(id)
	if err != nil {
		return err
	}
	return file.get(v)
}

//RemoveAll removes all files assosciated with this datastore
func (j *JSONFSStore) RemoveAll() error {
	if err := os.RemoveAll(j.dataDir); err != nil {
		return err
	}

	if err := os.MkdirAll(j.dataDir, 0755); err != nil {
		return err
	}
	return nil
}

//GetDataDir returns the directory path for files stored with this datastore
func (j *JSONFSStore) GetDataDir() string {
	return j.dataDir
}

func (j *JSONFSStore) AddToCollection(cPath string, id string, v interface{}) error {
	file, err := j.getJsonFile(cPath)
	if err != nil {
		return err
	}

	c, err := j.getCollection(file)
	if err != nil {
		return err
	}
	c[id] = v

	return file.set(c)
}

func (j *JSONFSStore) RemoveFromCollection(cPath string, id string) error {
	file, err := j.getJsonFile(cPath)
	if err != nil {
		return err
	}

	c, err := j.getCollection(file)
	if err != nil {
		return err
	}
	delete(c, id)

	return file.set(c)
}

func (j *JSONFSStore) GetCollection(cPath string) (map[string]interface{}, error) {
	file, err := j.getJsonFile(cPath)

	if err != nil {
		return nil, err
	}
	return j.getCollection(file)
}

func (j *JSONFSStore) getCollection(file *jsonFile) (map[string]interface{}, error) {
	c := make(map[string]interface{})

	if err := file.get(&c); err != nil {
		return nil, err
	}
	return c, nil
}

func (j *JSONFSStore) getJsonFile(path string) (*jsonFile, error) {
	if v, ok := j.files[path]; !ok {
		file, err := os.OpenFile(
			filepath.Join(j.dataDir, fmt.Sprintf("%s.json", path)),
			os.O_RDWR|os.O_CREATE,
			0644,
		)

		if err != nil {
			file.Close()
			debug.PrintStack()
			return nil, err
		}

		f := newJsonFile(file)

		j.files[path] = f

		return f, nil
	} else {
		return v, nil
	}

}

func (j *JSONFSStore) Flush() {
	for _, f := range j.files {
		f.flush()
	}
}

func (j *JSONFSStore) Stop() {
	for _, f := range j.files {
		f.stop()
	}
}
