package tests

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"net/url"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/nkanaev/yarr/src/storage"
)

func dbtest(t *testing.T, testcase func(t *testing.T, db storage.Storage)) {
	testurls := map[string]string{
		"sqlite": ":memory:",
	}

	if pgUrl := os.Getenv("YARR_POSTGRES_TEST_URL"); pgUrl != "" {
		dburl, cleanup, err := createPostgresDB(pgUrl)
		if err != nil {
			t.Fatalf("failed to create postgres test database: %v", err)
		}
		t.Cleanup(cleanup)
		testurls["postgres"] = dburl
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

func createPostgresDB(pgUrl string) (string, func(), error) {
	u, err := url.Parse(pgUrl)
	if err != nil {
		return "", nil, err
	}

	u.Path = "/postgres"
	adminConnStr := u.String()

	adminDB, err := sql.Open("postgres", adminConnStr)
	if err != nil {
		return "", nil, fmt.Errorf("admin connect: %w", err)
	}

	b := make([]byte, 4)
	if _, err := rand.Read(b); err != nil {
		adminDB.Close()
		return "", nil, fmt.Errorf("generate suffix: %w", err)
	}

	testDBName := "yarr_test_" + hex.EncodeToString(b)

	if _, err := adminDB.Exec(fmt.Sprintf(`CREATE DATABASE "%s"`, testDBName)); err != nil {
		adminDB.Close()
		return "", nil, fmt.Errorf("create database: %w", err)
	}
	adminDB.Close()

	u.Path = "/" + testDBName
	testURL := u.String()

	cleanup := func() {
		dropDB, err := sql.Open("postgres", adminConnStr)
		if err != nil {
			return
		}
		defer dropDB.Close()
		dropDB.Exec(fmt.Sprintf(
			`SELECT pg_terminate_backend(pid) FROM pg_stat_activity WHERE datname = '%s' AND pid <> pg_backend_pid()`,
			testDBName,
		))
		dropDB.Exec(fmt.Sprintf(`DROP DATABASE IF EXISTS "%s"`, testDBName))
	}

	return testURL, cleanup, nil
}
