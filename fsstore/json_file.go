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
	closed bool
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
	data, err := j.getBytes()
	if err != nil {
		return err
	}

	if data.Len() == 0 {
		return nil
	}
	decoder := json.NewDecoder(data)
	return decoder.Decode(v)
}

func (j *jsonFile) getBytes() (*bytes.Reader, error) {
	j.lock.RLock()
	defer j.lock.RUnlock()
	var data []byte
	var err error

	if j.isDirty() {
		data = j.buf.Bytes()
		log.Printf("[R] - buffer - %s", j.file.Name())
	} else {
		log.Printf("[R] - file - %s", j.file.Name())
		//if j.isEmpty() {
		//	return nil, nil
		//}
		j.file.Seek(0, 0)
		data, err = ioutil.ReadAll(j.file)
	}

	if err != nil {
		return nil, err
	}
	return bytes.NewReader(data), err
}

func (j *jsonFile) set(data interface{}) error {
	j.lock.Lock()
	defer j.lock.Unlock()

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	log.Printf("[W] %d bytes -> %s", len(jsonData), j.file.Name())
	j.buf.Reset()
	j.buf.Write(jsonData)
	return nil
}

func (j *jsonFile) remove() {
	j.stop()
	os.Remove(j.file.Name())
}

func (j *jsonFile) close() {
	j.file.Sync()
	j.file.Close()
	j.closed = true
}

func (j *jsonFile) isClosed() bool {
	return j.closed
}

func (j *jsonFile) start() {
	defer j.close()
	for range j.ticker.C {
		j.flush()
	}
}

func (j *jsonFile) stop() {
	j.ticker.Stop()
}

func (j *jsonFile) isDirty() bool {
	return j.buf.Len() > 0
}

func (j *jsonFile) isEmpty() bool {
	if j.isDirty() {
		return false
	}
	stat, _ := j.file.Stat()
	return stat.Size() == 0
}

func (j *jsonFile) flush() {
	j.lock.Lock()
	defer j.lock.Unlock()
	if j.isDirty() {
		log.Printf("[F] %d bytes -> %s ", j.buf.Len(), j.file.Name())
		j.file.Truncate(0)
		j.file.Seek(0, 0)
		j.buf.WriteTo(j.file)
		j.buf.Reset()
		j.stop()
	}

}
