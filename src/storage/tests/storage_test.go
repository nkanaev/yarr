package tests

import (
	"io"
	"log"
	"os"
	"testing"

	"github.com/nkanaev/yarr/src/storage"
)

func dbtest(t *testing.T, testcase func(t *testing.T, db storage.Storage)) {
	testurls := map[string]string {
		"sqlite": ":memory:",
		"postgres": "postgres://postgres:postgres@localhost:5432/yarr_test",
	}
	for testname, url := range testurls {
		db, err := storage.New(url)
		if err != nil {
			t.Fatalf("failed to init storage for %s: %v", url, err)
		}
		t.Run(testname, func(t *testing.T) {
			testcase(t, db)
		})
	}
}
