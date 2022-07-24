package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	telegramApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/reinanhs/golang-telegram-bot-example/internal/application/dto"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const (
	TestToken      = "153667468:AAHlSHlMqSt1f_uFmVRJbm5gntu2HI4WW8I"
	TestAppVersion = "1.0.1"
)

func createJsonToParseTelegramRequest(updateId int) *bytes.Reader {
	body, _ := json.Marshal(map[string]interface{}{
		"update_id": updateId,
	})

	return bytes.NewReader(body)
}

// TestParseTelegramRequest handles incoming update from the Telegram web hook
func TestParseTelegramRequest(t *testing.T) {
	// Table driven tests
	var tests = []struct {
		name    string
		request *http.Request
		want    *telegramApi.Update
		err     error
	}{
		{
			"empty body",
			httptest.NewRequest(http.MethodPost, "/", nil),
			(*telegramApi.Update)(nil),
			errors.New("EOF"),
		},
		{
			"invalid json character",
			httptest.NewRequest(http.MethodPost, "/", strings.NewReader("a=1")),
			(*telegramApi.Update)(nil),
			errors.New("invalid character 'a' looking for beginning of value"),
		},
		{
			"invalid http get method",
			httptest.NewRequest(http.MethodGet, "/", nil),
			(*telegramApi.Update)(nil),
			errors.New("unsupported method " + http.MethodGet),
		},
		{
			"valid json with update id",
			httptest.NewRequest(http.MethodPost, "/", createJsonToParseTelegramRequest(1)),
			&telegramApi.Update{UpdateID: 1},
			nil,
		},
		{
			"invalid update id",
			httptest.NewRequest(http.MethodPost, "/", createJsonToParseTelegramRequest(0)),
			nil,
			errors.New("invalid update id of 0 indicates failure to parse incoming update"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseTelegramRequest(tt.request)

			if err != nil && assert.Error(t, err) {
				assert.Equal(t, err.Error(), tt.err.Error(), "they should be equal")
			}

			if tt.want != got {
				assert.Equal(t, got, tt.want, "they should be equal")
			}
		})
	}
}

// TestResponseMessage responsible for returning a message about the status of the request
func TestResponseMessage(t *testing.T) {
	// Table driven tests
	var tests = []struct {
		name    string
		version string
		message string
		status  int
	}{
		{
			name:    "check status 200",
			version: TestAppVersion,
			status:  http.StatusOK,
			message: http.StatusText(http.StatusOK),
		},
		{
			name:    "check status 422",
			version: TestAppVersion,
			status:  http.StatusUnprocessableEntity,
			message: http.StatusText(http.StatusUnprocessableEntity),
		},
		{
			name:    "check empty version",
			version: "",
			status:  http.StatusOK,
			message: http.StatusText(http.StatusOK),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("APP_VERSION", tt.version)

			w := httptest.NewRecorder()
			responseMessage(tt.message, tt.status, w)

			resp := w.Result()
			body, _ := io.ReadAll(resp.Body)

			if tt.version == "" {
				tt.version = "1.0.0"
			}

			respData := dto.Response{
				Version: tt.version,
				Data: dto.Data{
					Data: dto.ResponseMessage{
						Message: tt.message,
						Status:  tt.status,
					},
				},
			}

			expectedJson, _ := json.Marshal(respData)
			expectedStatus := fmt.Sprintf("%d %s", tt.status, http.StatusText(tt.status))

			assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
			assert.Equal(t, expectedStatus, resp.Status)
			assert.Equal(t, string(expectedJson), string(body))
		})
	}
}

func TestHandleTelegramWebHook(t *testing.T) {
	// Table driven tests
	var tests = []struct {
		name          string
		message       string
		status        int
		telegramToken string
		body          io.Reader
	}{
		{
			name:          "check with valid body",
			message:       "Operation performed successfully",
			status:        http.StatusOK,
			telegramToken: TestToken,
			body:          createJsonToParseTelegramRequest(1),
		},
		{
			name:          "invalid token",
			message:       "Could not start bot",
			status:        http.StatusInternalServerError,
			telegramToken: "",
			body:          createJsonToParseTelegramRequest(1),
		},
		{
			name:          "invalid body",
			message:       "Could not read update request",
			status:        http.StatusUnprocessableEntity,
			telegramToken: TestToken,
			body:          strings.NewReader("a=1"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("APP_VERSION", TestAppVersion)
			t.Setenv("TELEGRAM_BOT_TOKEN", tt.telegramToken)

			req := httptest.NewRequest(http.MethodPost, "/", tt.body)
			w := httptest.NewRecorder()

			HandleTelegramWebHook(w, req)

			resp := w.Result()
			body, _ := io.ReadAll(resp.Body)

			respData := dto.Response{
				Version: TestAppVersion,
				Data: dto.Data{
					Data: dto.ResponseMessage{
						Message: tt.message,
						Status:  tt.status,
					},
				},
			}
			expectedJson, _ := json.Marshal(respData)

			assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
			assert.Equal(t, tt.status, resp.StatusCode)
			assert.Equal(t, string(expectedJson), string(body))
		})
	}
}
