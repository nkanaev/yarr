package storage

import (
	"io"
	"log"
	"os"
	"testing"
)

func testDB() *Storage {
	log.SetOutput(io.Discard)
	db, _ := New(":memory:")
	log.SetOutput(os.Stderr)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	return db
}

func TestStorage(t *testing.T) {
	db, err := New(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	if db == nil {
		t.Fatal("no db")
	}
}
