package app

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestApp_Start(t *testing.T) {
	ctx := context.Background()

	valkey, postgres, backend := initApp(t, ctx)

	app := New()
	app.Start()

	t.Cleanup(func() {
		app.Stop()
		if err := backend.Terminate(ctx); err != nil {
			panic(err)
		}
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
		response: fmt.Sprintf(`{"expression":{"id":"%s","expression":"2+2","result":4,"status":"DONE"}}`, idDone),
	})

	sendRequest(t, httpClient, test{
		name:     "calc_FAILED",
		url:      fmt.Sprintf("%s/%s", expressionsUrl, idFailed),
		method:   http.MethodGet,
		auth:     header,
		status:   http.StatusOK,
		response: fmt.Sprintf(`{"expression":{"id":"%s","expression":"2/(2-2)","result":0,"status":"FAILED"}}`, idFailed),
	})
}
