package app

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"net/http"
	"orchestrator/internal/config"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

func TestGrpcApp(t *testing.T) {
	ctx := context.Background()

	_, filename, _, _ := runtime.Caller(0)
	rootDir := filepath.Join(filepath.Dir(filename), "../..")
	err := os.Chdir(rootDir)
	require.NoError(t, err)

	valkey, postgres := createDb(t, ctx)

	cfg := config.New()
	app := New(cfg)

	app.Start()

	agentContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image: "linuxfight/yndxfinal-agent:main",
			HostConfigModifier: func(hc *container.HostConfig) {
				hc.ExtraHosts = []string{"host.docker.internal:host-gateway"}
			},
			Env: map[string]string{
				"API_ADDR": "host.docker.internal:9090",
			},
			WaitingFor: wait.ForLog("Worker started with URL: host.docker.internal:9090").WithStartupTimeout(30 * time.Second),
		},
		Started: true,
	})
	require.NoError(t, err)

	t.Cleanup(func() {
		if err := agentContainer.Terminate(ctx); err != nil {
			panic(err)
		}
		app.Stop()
		if err := postgres.Terminate(ctx); err != nil {
			panic(err)
		}
		if err := valkey.Terminate(ctx); err != nil {
			panic(err)
		}
	})

	httpClient := &http.Client{
		Transport: &http.Transport{
			DisableCompression: true,
		},
	}

	token := getToken(t, httpClient)
	header := fmt.Sprintf("Bearer %s", token)
	idDone := getTaskId(t, httpClient, header, `{"expression": "2+2"}`)
	idFailed := getTaskId(t, httpClient, header, `{"expression": "2/(2-2)"}`)

	time.Sleep(10 * time.Second)

	sendRequest(t, httpClient, test{
		name:     "calc_DONE",
		url:      fmt.Sprintf("%s/%s", expressionsUrl, idDone),
		method:   http.MethodGet,
		auth:     header,
		status:   http.StatusOK,
		response: fmt.Sprintf(`{"expression":{"id":"%s","result":4,"status":"DONE"}}`, idDone),
	})

	sendRequest(t, httpClient, test{
		name:     "calc_FAILED",
		url:      fmt.Sprintf("%s/%s", expressionsUrl, idFailed),
		method:   http.MethodGet,
		auth:     header,
		status:   http.StatusOK,
		response: fmt.Sprintf(`{"expression":{"id":"%s","result":0,"status":"FAILED"}}`, idFailed),
	})
}
