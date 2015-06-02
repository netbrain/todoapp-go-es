package fsstore

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"code.google.com/p/go-uuid/uuid"
)

func TestCanCreateAndGetAndRemove(t *testing.T) {
	t.Parallel()
	dataDir, _ := ioutil.TempDir(os.TempDir(), "")
	fs, err := NewJSONFSStore(dataDir)
	if err != nil {
		t.Fatal(err)
	}
	id := uuid.New()
	err = fs.Set(id, struct{ A string }{A: "Test"})
	if err != nil {
		t.Fatal(err)
	}

	fs.Flush()
	stat, _ := os.Stat(filepath.Join(dataDir, fmt.Sprintf("%s.json", id)))
	if stat.Size() == 0 {
		t.Fatal("File not stored")
	}

	defer func() {
		if err := fs.Remove(id); err != nil {
			t.Fatal(err)
		}

		if _, err := os.Stat(filepath.Join(dataDir, id)); !os.IsNotExist(err) {
			t.Fatal("Removed file exists..")
		}

		if err := fs.RemoveAll(); err != nil {
			t.Fatal(err)
		}

		if _, err := os.Stat(dataDir); !os.IsNotExist(err) {
			t.Fatal("Removed dir exists..")
		}
	}()

	v := &struct{ A string }{}
	if err := fs.Get(id, v); err != nil {
		t.Fatal(err)
	}

	if v.A != "Test" {
		t.Fatal("Stored struct not equal to expected")
	}

}

func TestCollectionAddAndRemove(t *testing.T) {
	t.Parallel()
	dataDir, _ := ioutil.TempDir(os.TempDir(), "")
	fs, err := NewJSONFSStore(dataDir)
	if err != nil {
		t.Fatal(err)
	}

	id := uuid.New()
	collectionName := "test-collection"
	err = fs.AddToCollection(collectionName, id, struct{ A string }{A: "Test"})
	if err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(filepath.Join(dataDir, fmt.Sprintf("%s.json", collectionName))); os.IsNotExist(err) {
		t.Fatal(err)
	}

	fs.RemoveFromCollection(collectionName, id)
	fs.Flush()

	if stat, err := os.Stat(filepath.Join(dataDir, fmt.Sprintf("%s.json", collectionName))); os.IsNotExist(err) {
		t.Fatal(err)
	} else if stat.Size() > 2 {
		t.Fatalf("Unexpected file size: %d", stat.Size())
	}

}

func TestAddSeveralToCollection(t *testing.T) {
	t.Parallel()
	dataDir, _ := ioutil.TempDir(os.TempDir(), "")
	fs, err := NewJSONFSStore(dataDir)
	if err != nil {
		t.Fatal(err)
	}

	collectionName := "test-collection"

	for x := 0; x < 10; x++ {
		id := uuid.New()
		err := fs.AddToCollection(collectionName, id, struct{ A string }{A: "Test"})
		if err != nil {
			t.Fatal(err)
		}
	}

	m, err := fs.GetCollection(collectionName)
	if err != nil {
		t.Fatal(err)
	}
	if len(m) != 10 {
		t.Fatalf("Expected 10 elements in collection, but got: %d", len(m))
	}

}
