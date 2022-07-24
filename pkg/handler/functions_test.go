package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	telegramApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
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