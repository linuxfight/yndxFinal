package app

import (
	"context"
	"github.com/bytedance/sonic"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/network"
	"github.com/testcontainers/testcontainers-go/wait"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

type test struct {
	name     string
	url      string
	method   string
	body     string
	auth     string
	status   int
	response string
}

type authResponse struct {
	Token string `json:"token"`
}

type calculateResponse struct {
	Id string `json:"id"`
}

const (
	registerUrl    = "http://localhost:8080/api/v1/register"
	calculateUrl   = "http://localhost:8080/api/v1/calculate"
	expressionsUrl = "http://localhost:8080/api/v1/expressions"
)

func sendRequest(t *testing.T, client *http.Client, testCase test) {
	t.Run(testCase.name, func(t *testing.T) {
		var req *http.Request
		var err error
		switch testCase.method {
		case http.MethodPost:
			req, err = http.NewRequest(testCase.method, testCase.url, strings.NewReader(testCase.body))
		case http.MethodGet:
			req, err = http.NewRequest(testCase.method, testCase.url, nil)
		default:
			t.Fatalf("invalid method: %s", testCase.method)
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

		if testCase.response != "" {
			bytes, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			if testCase.response != string(bytes) {
				t.Errorf("expected response %s but got %s", testCase.response, string(bytes))
			}
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

	var body authResponse
	err = sonic.Unmarshal(bytes, &body)
	require.NoError(t, err)

	return body.Token
}

func getTaskId(t *testing.T, client *http.Client, header, expression string) string {
	req, err := http.NewRequest(http.MethodPost, calculateUrl, strings.NewReader(expression))
	require.NoError(t, err)
	req.Header.Set("Authorization", header)
	resp, err := client.Do(req)
	require.NoError(t, err)

	bytes, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected status code %d but got %d", http.StatusCreated, resp.StatusCode)
	}

	var body calculateResponse
	err = sonic.Unmarshal(bytes, &body)
	require.NoError(t, err)

	return body.Id
}

func initApp(t *testing.T, ctx context.Context) (testcontainers.Container, testcontainers.Container, testcontainers.Container) {
	appNet, err := network.New(ctx)
	require.NoError(t, err)

	valkeyContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "valkey/valkey:latest",
			ExposedPorts: []string{"6379/tcp"},
			Name:         "valkey",
			Networks: []string{
				appNet.Name,
			},
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
			Name:         "postgres",
			Env: map[string]string{
				"POSTGRES_DB":       "db",
				"POSTGRES_USER":     "postgres",
				"POSTGRES_PASSWORD": "password",
			},
			Networks: []string{
				appNet.Name,
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

	backend, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image: "linuxfight/yndxfinal-orchestrator:main",
			Networks: []string{
				appNet.Name,
			},
			Env: map[string]string{
				"POSTGRES_CONN": "postgres://postgres:password@postgres:5432/db",
				"VALKEY_CONN":   "valkey:6379",
			},
			ExposedPorts: []string{"8080/tcp", "9090/tcp"},
			HostConfigModifier: func(hc *container.HostConfig) {
				hc.PortBindings = nat.PortMap{
					"8080/tcp": []nat.PortBinding{
						{
							HostIP:   "127.0.0.1",
							HostPort: "8080",
						},
					},
					"9090/tcp": []nat.PortBinding{
						{
							HostIP:   "127.0.0.1",
							HostPort: "9090",
						},
					},
				}
			},
			WaitingFor: wait.ForHTTP("/startupz").WithPort("8080").WithStartupTimeout(30 * time.Second),
		},
		Started: true,
	})
	require.NoError(t, err)

	return valkeyContainer, postgresContainer, backend
}
