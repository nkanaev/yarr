package storage

import (
	"io"
	"log"
	"os"
	"testing"
)

func testDB() *Storage {
	log.SetOutput(io.Discard)
	db, err := New(":memory:")
	if err != nil {
		panic(err)
	}
	log.SetOutput(os.Stderr)
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
