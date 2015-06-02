package fsstore

import (
	"io/ioutil"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestDirtyness(t *testing.T) {
	t.Parallel()

	tempFile, _ := ioutil.TempFile(os.TempDir(), "")
	f := newJsonFile(tempFile)
	if f.isDirty() {
		t.Fatal("should not be dirty")
	}

	f.set(struct{}{})

	// TODO this concept of dirty needs to go
	//if !f.isDirty() {
	//	t.Fatal("should be dirty")
	//}
}

func TestEmptyness(t *testing.T) {
	t.Parallel()

	tempFile, _ := ioutil.TempFile(os.TempDir(), "")
	f := newJsonFile(tempFile)

	if !f.isEmpty() {
		t.Fatal("should be empty")
	}

	f.set(struct{}{})

	// Dirty handling
	//	if !f.isEmpty() {
	//		t.Fatal("should still be empty")
	//	}

	//	f.flush()

	if f.isEmpty() {
		t.Fatal("should not be empty")
	}
}

func TestCanGetWrittenData(t *testing.T) {
	t.Parallel()

	tempFile, _ := ioutil.TempFile(os.TempDir(), "")
	f := newJsonFile(tempFile)

	s1 := struct{ A bool }{A: true}
	var s2 struct{ A bool }

	f.set(s1)

	f.get(&s2)

	if !s2.A {
		t.Fatal("A is not true")
	}

}

func TestCanGetWrittenCollection(t *testing.T) {
	t.Parallel()

	tempFile, _ := ioutil.TempFile(os.TempDir(), "")
	f := newJsonFile(tempFile)

	c1 := make(map[string]interface{})
	c2 := make(map[string]interface{})
	s1 := struct{ A bool }{A: true}

	for x := 0; x < 10; x++ {
		c1[strconv.Itoa(x)] = s1
		f.set(c1)
	}

	time.Sleep(100)

	f.get(&c2)

	if len(c2) != 10 {
		t.Fatalf("expected 10 results, got: %d", len(c2))
	}

}

func BenchmarkWrite(b *testing.B) {
	tempFile, _ := ioutil.TempFile(os.TempDir(), "")
	f := newJsonFile(tempFile)

	for n := 0; n < b.N; n++ {
		f.set(struct{ A string }{A: "test"})
		f.flush()
	}

}
