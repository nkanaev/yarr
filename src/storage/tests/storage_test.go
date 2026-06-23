package tests

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/nkanaev/yarr/src/storage"
)

func dbtest(t *testing.T, testcase func(t *testing.T, db storage.Storage)) {
	t.Parallel()
	testurls := map[string]string{
		"sqlite": ":memory:",
	}

	if pgImage := os.Getenv("YARR_POSTGRES_TEST_IMAGE"); pgImage != "" {
		dburl, cleanup := startPostgresContainer(t, pgImage)
		t.Cleanup(cleanup)
		testurls["postgres"] = dburl
	} else if !testing.Short() {
		t.Fatalf("YARR_POSTGRES_TEST_IMAGE not set; use -short to skip docker tests")
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

func startPostgresContainer(t *testing.T, image string) (string, func()) {
	// database credentials
	dbUser := "testuser"
	dbPass := "password"
	dbName := "yarrtest"

	// generate unique container name
	testHash := sha256.Sum256([]byte(t.Name()))
	containerName := fmt.Sprintf("yarr-test-pg-%x-%d", testHash[:8], time.Now().UnixNano())

	cmd := exec.Command(
		"docker", "run", "-d", "--rm",
		"--name", containerName,
		"-p", "0:5432",
		"-e", "POSTGRES_USER="+dbUser,
		"-e", "POSTGRES_PASSWORD="+dbPass,
		"-e", "POSTGRES_DB="+dbName,
		image,
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("failed to start postgres container: %v\n%s", err, string(out))
	}

	// retrieve the host port assigned by docker
	portCmd := exec.Command("docker", "port", containerName, "5432/tcp")
	portOut, err := portCmd.Output()
	if err != nil {
		t.Fatalf("failed to get container port: %v", err)
	}
	parts := strings.Split(strings.TrimSpace(string(portOut)), ":")
	dbPort := parts[len(parts)-1]

	// build connection string
	pgUrl := fmt.Sprintf(
		"postgres://%s:%s@localhost:%s/%s?sslmode=disable",
		dbUser,
		dbPass,
		dbPort,
		dbName,
	)

	// wait up to 15 seconds for the container to accept connections
	deadline := time.Now().Add(15 * time.Second)
	for time.Now().Before(deadline) {
		db, err := sql.Open("postgres", pgUrl)
		if err != nil {
			continue
		}
		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		err = db.PingContext(ctx)
		cancel()
		db.Close()
		if err == nil {
			goto ready
		}
		time.Sleep(200 * time.Millisecond)
	}
	t.Fatalf("timed out waiting for postgres container to be ready")

ready:
	// return connection url and a cleanup function that stops the container
	return pgUrl, func() {
		stop := exec.Command("docker", "stop", containerName)
		if err := stop.Run(); err != nil {
			t.Logf("failed to stop container %s: %v", containerName, err)
		}
	}
}
