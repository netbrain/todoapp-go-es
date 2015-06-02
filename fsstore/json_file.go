package fsstore

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"
)

type jsonFile struct {
	file   *os.File
	buf    *bytes.Buffer
	lock   *sync.RWMutex
	ticker *time.Ticker
}

func newJsonFile(file *os.File) *jsonFile {
	f := &jsonFile{
		file:   file,
		buf:    new(bytes.Buffer),
		lock:   new(sync.RWMutex),
		ticker: time.NewTicker(time.Second * 5),
	}
	go f.start()
	return f
}

func (j *jsonFile) get(v interface{}) error {
	j.lock.RLock()
	defer j.lock.RUnlock()
	var data []byte
	var err error

	log.Printf("[R] %s", j.file.Name())
	if j.isDirty() {
		data, err = ioutil.ReadAll(j.buf)
	} else {
		if j.isEmpty() {
			return nil
		}
		j.file.Seek(0, 0)
		data, err = ioutil.ReadAll(j.file)
	}

	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

func (j *jsonFile) set(data interface{}) error {
	jsonData, err := json.Marshal(data)
	log.Printf("[W] %d bytes -> %s", len(jsonData), j.file.Name())
	if err != nil {
		return err
	}

	j.lock.Lock()

	j.buf.Reset()
	j.buf.Write(jsonData)
	j.lock.Unlock()
	j.flush()
	return nil
}

func (j *jsonFile) remove() {
	//	j.writeChan <- nil
	//	j.stop()
}

func (j *jsonFile) close() {
	j.file.Sync()
	j.file.Close()
}

func (j *jsonFile) start() {
	for range j.ticker.C {
		j.flush()
	}
	j.close()
}

func (j *jsonFile) stop() {
	j.ticker.Stop()
}

func (j *jsonFile) isDirty() bool {
	return j.buf.Len() > 0
}

func (j *jsonFile) isEmpty() bool {
	stat, _ := j.file.Stat()
	return stat.Size() == 0
}

func (j *jsonFile) flush() {
	if j.isDirty() {
		j.lock.Lock()
		defer j.lock.Unlock()
		log.Printf("[F] %d bytes -> %s ", j.buf.Len(), j.file.Name())
		j.file.Truncate(0)
		j.file.Seek(0, 0)
		j.buf.WriteTo(j.file)
		j.buf.Reset()
	}

}
