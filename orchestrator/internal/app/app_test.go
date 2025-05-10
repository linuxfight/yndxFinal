package app

import (
	"context"
	"fmt"
	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/require"
	"net/http"
	"orchestrator/internal/config"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

type test struct {
	name   string
	url    string
	method string
	body   string
	auth   string
	status int
}

const (
	loginUrl       = "http://localhost:8080/api/v1/login"
	registerUrl    = "http://localhost:8080/api/v1/register"
	calculateUrl   = "http://localhost:8080/api/v1/calculate"
	expressionsUrl = "http://localhost:8080/api/v1/expressions"
)

func TestApp(t *testing.T) {
	ctx := context.Background()

	_, filename, _, _ := runtime.Caller(0)
	rootDir := filepath.Join(filepath.Dir(filename), "../..")
	err := os.Chdir(rootDir)
	require.NoError(t, err)

	valkey, postgres := createDb(t, ctx)

	cfg := config.New()
	app := New(cfg)

	httpClient := &http.Client{
		Transport: &http.Transport{
			DisableCompression: true,
		},
	}

	t.Cleanup(func() {
		app.Stop()
		if err := postgres.Terminate(ctx); err != nil {
			panic(err)
		}
		if err := valkey.Terminate(ctx); err != nil {
			panic(err)
		}
	})

	app.Start()

	tests := generateTests(t, httpClient)

	for _, testCase := range tests {
		sendRequest(t, httpClient, testCase)
	}
}

func generateTests(t *testing.T, client *http.Client) []test {
	authTests := []test{
		// register
		// @Success      201
		// @Failure      400
		// @Failure      409
		// @Failure      422
		{
			name:   "register_200",
			url:    registerUrl,
			method: http.MethodPost,
			body:   `{"login": "login", "password": "password"}`,
			auth:   "",
			status: http.StatusCreated,
		},
		{
			name:   "register_400",
			url:    registerUrl,
			method: http.MethodPost,
			body:   `{"login": "", "password": ""}`,
			auth:   "",
			status: http.StatusBadRequest,
		},
		{
			name:   "register_409",
			url:    registerUrl,
			method: http.MethodPost,
			body:   `{"login": "login", "password": "password"}`,
			auth:   "",
			status: http.StatusConflict,
		},
		{
			name:   "register_422",
			url:    registerUrl,
			method: http.MethodPost,
			body:   `{"login": "", "password": "`,
			auth:   "",
			status: http.StatusUnprocessableEntity,
		},
		// login
		// @Success      200
		// @Failure      400
		// @Failure      401
		// @Failure      404
		// @Failure      422
		{
			name:   "login_200",
			url:    loginUrl,
			method: http.MethodPost,
			body:   `{"login": "login", "password": "password"}`,
			auth:   "",
			status: http.StatusOK,
		},
		{
			name:   "login_400",
			url:    loginUrl,
			method: http.MethodPost,
			body:   `{"login": "", "password": ""}`,
			auth:   "",
			status: http.StatusBadRequest,
		},
		{
			name:   "login_401",
			url:    loginUrl,
			method: http.MethodPost,
			body:   `{"login": "login", "password": "password1231321"}`,
			auth:   "",
			status: http.StatusUnauthorized,
		},
		{
			name:   "login_404",
			url:    loginUrl,
			method: http.MethodPost,
			body:   `{"login": "login321321", "password": "password1231321"}`,
			auth:   "",
			status: http.StatusNotFound,
		},
		{
			name:   "login_422",
			url:    loginUrl,
			method: http.MethodPost,
			body:   `{"login": "", "password": "`,
			auth:   "",
			status: http.StatusUnprocessableEntity,
		},
	}

	token := getToken(t, client)
	header := fmt.Sprintf("Bearer %s", token)
	id := getTaskId(t, client, header)

	calcTests := []test{
		// calc
		// @Success      200
		// @Success      201
		// @Failure      400
		// @Failure      403
		// @Failure      422
		{
			name:   "calc_201",
			url:    calculateUrl,
			method: http.MethodPost,
			body:   `{"expression": "2+2"}`,
			auth:   header,
			status: http.StatusCreated,
		},
		{
			name:   "calc_200",
			url:    calculateUrl,
			method: http.MethodPost,
			body:   `{"expression": "2+2"}`,
			auth:   header,
			status: http.StatusOK,
		},
		{
			name:   "calc_400",
			url:    calculateUrl,
			method: http.MethodPost,
			body:   `{"expression": "2+2dsadsadsa"}`,
			auth:   header,
			status: http.StatusBadRequest,
		},
		{
			name:   "calc_403",
			url:    calculateUrl,
			method: http.MethodPost,
			body:   `{"expression": "2+2"}`,
			auth:   "",
			status: http.StatusForbidden,
		},
		{
			name:   "calc_422",
			url:    calculateUrl,
			method: http.MethodPost,
			body:   `{"expression": "2+2`,
			auth:   header,
			status: http.StatusUnprocessableEntity,
		},
		// list
		// @Success      200
		// @Success      403
		{
			name:   "list_200",
			url:    expressionsUrl,
			method: http.MethodGet,
			auth:   header,
			status: http.StatusOK,
		},
		{
			name:   "list_403",
			url:    expressionsUrl,
			method: http.MethodGet,
			auth:   "",
			status: http.StatusForbidden,
		},
		// id
		// @Success      200
		// @Failure      400
		// @Failure      403
		// @Failure      404
		{
			name:   "id_200",
			url:    fmt.Sprintf("%s/%s", expressionsUrl, id),
			method: http.MethodGet,
			auth:   header,
			status: http.StatusOK,
		},
		{
			name:   "id_400_len",
			url:    fmt.Sprintf("%s/%s", expressionsUrl, "dsadsa"),
			method: http.MethodGet,
			auth:   header,
			status: http.StatusBadRequest,
		},
		{
			name:   "id_400_inv_id",
			url:    fmt.Sprintf("%s/%s", expressionsUrl, "01^^VM4KDA0BQ4GAGY44BXQXMN"),
			method: http.MethodGet,
			auth:   header,
			status: http.StatusBadRequest,
		},
		{
			name:   "id_403",
			url:    fmt.Sprintf("%s/%s", expressionsUrl, id),
			method: http.MethodGet,
			auth:   "",
			status: http.StatusForbidden,
		},
		{
			name:   "id_404",
			url:    fmt.Sprintf("%s/%s", expressionsUrl, ulid.Make().String()),
			method: http.MethodGet,
			auth:   header,
			status: http.StatusNotFound,
		},
	}

	tests := append(authTests, calcTests...)

	return tests
}
