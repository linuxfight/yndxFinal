package app

import (
	"context"
	"github.com/bytedance/sonic"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"io"
	"net/http"
	"orchestrator/internal/controllers/dto"
	"strings"
	"testing"
	"time"
)

func sendRequest(t *testing.T, client *http.Client, testCase test) {
	t.Run(testCase.name, func(t *testing.T) {
		var req *http.Request
		var err error
		if testCase.method == http.MethodPost {
			req, err = http.NewRequest(testCase.method, testCase.url, strings.NewReader(testCase.body))
		} else if testCase.method == http.MethodGet {
			req, err = http.NewRequest(testCase.method, testCase.url, nil)
		}
		require.NoError(t, err)

		if testCase.auth != "" {
			req.Header.Set("Authorization", testCase.auth)
		}

		resp, err := client.Do(req)
		require.NoError(t, err)

		if resp.StatusCode != testCase.status {
			t.Errorf("expected status code %d but got %d", testCase.status, resp.StatusCode)
		}
	})
}

func getToken(t *testing.T, client *http.Client) string {
	loginBody := `{"login": "calc", "password": "calc"}`
	req, err := http.NewRequest(http.MethodPost, registerUrl, strings.NewReader(loginBody))
	require.NoError(t, err)
	resp, err := client.Do(req)
	require.NoError(t, err)

	bytes, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected status code %d but got %d", http.StatusCreated, resp.StatusCode)
	}

	var body dto.AuthResponse
	err = sonic.Unmarshal(bytes, &body)
	require.NoError(t, err)

	return body.Token
}

func getTaskId(t *testing.T, client *http.Client, header string) string {
	calcBody := `{"expression": "52+52"}`

	req, err := http.NewRequest(http.MethodPost, calculateUrl, strings.NewReader(calcBody))
	require.NoError(t, err)
	req.Header.Set("Authorization", header)
	resp, err := client.Do(req)
	require.NoError(t, err)

	bytes, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected status code %d but got %d", http.StatusCreated, resp.StatusCode)
	}

	var body dto.CalculateResponse
	err = sonic.Unmarshal(bytes, &body)
	require.NoError(t, err)

	return body.Id
}

func createDb(t *testing.T, ctx context.Context) (testcontainers.Container, testcontainers.Container) {
	valkeyContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "valkey/valkey:latest",
			ExposedPorts: []string{"6379/tcp"},
			HostConfigModifier: func(hc *container.HostConfig) {
				hc.PortBindings = nat.PortMap{
					"6379/tcp": []nat.PortBinding{
						{
							HostIP:   "127.0.0.1",
							HostPort: "6379",
						},
					},
				}
			},
			WaitingFor: wait.ForLog("Ready to accept connections").WithStartupTimeout(30 * time.Second),
		},
		Started: true,
	})
	require.NoError(t, err)

	postgresContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres:alpine",
			ExposedPorts: []string{"5432/tcp"},
			Env: map[string]string{
				"POSTGRES_DB":       "db",
				"POSTGRES_USER":     "postgres",
				"POSTGRES_PASSWORD": "password",
			},
			WaitingFor: wait.ForAll(
				wait.ForLog("database system is ready to accept connections"),
				wait.ForListeningPort("5432/tcp"),
			).WithDeadline(30 * time.Second),
			HostConfigModifier: func(hc *container.HostConfig) {
				hc.PortBindings = nat.PortMap{
					"5432/tcp": []nat.PortBinding{
						{
							HostIP:   "127.0.0.1",
							HostPort: "5432",
						},
					},
				}
			},
		},
		Started: true,
	})
	require.NoError(t, err)

	return valkeyContainer, postgresContainer
}
