package utils

import (
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/require"
)

func TestSendError(t *testing.T) {
	tests := []struct {
		name        string
		status      int
		message     string
		wantCode    int
		wantMessage string
	}{
		{
			name:        "standard error",
			status:      fiber.StatusBadRequest,
			message:     "Invalid input",
			wantCode:    fiber.StatusBadRequest,
			wantMessage: "Invalid input",
		},
		{
			name:        "not found error",
			status:      fiber.StatusNotFound,
			message:     "Resource not found",
			wantCode:    fiber.StatusNotFound,
			wantMessage: "Resource not found",
		},
		{
			name:        "internal server error",
			status:      fiber.StatusInternalServerError,
			message:     "Something went wrong",
			wantCode:    fiber.StatusInternalServerError,
			wantMessage: "Something went wrong",
		},
		{
			name:        "empty message",
			status:      fiber.StatusBadRequest,
			message:     "",
			wantCode:    fiber.StatusBadRequest,
			wantMessage: "",
		},
		{
			name:        "custom status code",
			status:      418, // I'm a teapot
			message:     "Short and stout",
			wantCode:    418,
			wantMessage: "Short and stout",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create Fiber app
			app := fiber.New()

			// Setup test route
			app.Get("/test", func(c fiber.Ctx) error {
				return SendError(c, tt.message, tt.status)
			})

			// Create test request
			req := httptest.NewRequest("GET", "/test", nil)
			resp, err := app.Test(req)
			require.NoError(t, err)
			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {
					panic(err)
				}
			}(resp.Body)

			// Verify status code
			require.Equal(t, tt.wantCode, resp.StatusCode)

			// Verify content type
			require.Equal(t, fiber.MIMEApplicationJSON, resp.Header.Get(fiber.HeaderContentType))

			// Parse response body
			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			// Verify JSON structure
			var result fiber.Error
			err = app.Config().JSONDecoder(body, &result)
			require.NoError(t, err)

			// Verify error details
			require.Equal(t, tt.wantCode, result.Code)
			require.Equal(t, tt.wantMessage, result.Message)
		})
	}
}
