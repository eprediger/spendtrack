package features_test

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/cucumber/godog"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

var (
	apiURL string
	pool   *pgxpool.Pool
)

func TestFeatures(t *testing.T) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/spendtrack?sslmode=disable"
	}

	var err error
	pool, err = pgxpool.New(context.Background(), dbURL)
	require.NoError(t, err)
	require.NoError(t, pool.Ping(context.Background()))
	defer pool.Close()

	var stop func()
	apiURL, stop = startAPI(t, dbURL)
	defer stop()

	suite := godog.TestSuite{
		ScenarioInitializer: initializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"."},
			Tags:     "~@todo",
			TestingT: t,
		},
	}
	if suite.Run() != 0 {
		t.Fatal("godog suite failed")
	}
}

// startAPI builds and runs the real binary, then waits for it to
// answer the health endpoint. The suite exercises the HTTP surface
// exactly as a deployed client would.
func startAPI(t *testing.T, dbURL string) (string, func()) {
	t.Helper()

	bin := filepath.Join(t.TempDir(), "api")
	out, err := exec.Command("go", "build", "-o", bin, "../cmd/api").CombinedOutput()
	require.NoError(t, err, "api build failed: %s", out)

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	port := ln.Addr().(*net.TCPAddr).Port
	_ = ln.Close()

	cmd := exec.Command(bin)
	cmd.Env = append(os.Environ(),
		"DATABASE_URL="+dbURL,
		fmt.Sprintf("PORT=%d", port))
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	require.NoError(t, cmd.Start())

	base := fmt.Sprintf("http://127.0.0.1:%d", port)
	deadline := time.Now().Add(15 * time.Second)
	for {
		resp, err := http.Get(base + "/api/v1/health")
		if err == nil {
			_ = resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				return base, func() {
					_ = cmd.Process.Kill()
					_ = cmd.Wait()
				}
			}
		}
		if time.Now().After(deadline) {
			_ = cmd.Process.Kill()
			t.Fatal("api did not become healthy")
		}
		time.Sleep(100 * time.Millisecond)
	}
}
