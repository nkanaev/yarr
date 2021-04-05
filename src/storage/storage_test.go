package storage

import (
	"io"
	"log"
	"os"
)

func testDB() *Storage {
	log.SetOutput(io.Discard)
	db, err := New(":memory:")
	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Exit(1)
	}
	log.SetOutput(os.Stderr)
	return db
}
