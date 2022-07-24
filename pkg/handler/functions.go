package handler

import (
	"encoding/json"
	"errors"
	telegramApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/reinanhs/golang-telegram-bot-example/internal/application/dto"
	"log"
	"net/http"
	"os"
)

const telegramTokenEnv string = "TELEGRAM_BOT_TOKEN"

// HandleTelegramWebHook sends a message back to the chat with a punchline starting by the message provided by the user.
func HandleTelegramWebHook(w http.ResponseWriter, r *http.Request) {
	var bot, errBot = telegramApi.NewBotAPI(os.Getenv(telegramTokenEnv))
	if errBot != nil {
		w = responseMessage("Could not start bot", http.StatusInternalServerError, w)
		return
	}

	// Parse incoming request
	var update, err = parseTelegramRequest(r)
	if err != nil {
		w = responseMessage("Could not read update request", http.StatusUnprocessableEntity, w)
		return
	}

	ActionHandler(update, bot)
	w = responseMessage("Operation performed successfully", http.StatusOK, w)
	return
}

// responseMessage responsible for returning a message about the status of the request
func responseMessage(m string, s int, w http.ResponseWriter) http.ResponseWriter {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(s)

	version := os.Getenv("APP_VERSION")
	if version == "" {
		version = "1.0.0"
	}

	message := dto.ResponseMessage{
		Message: m,
		Status:  s,
	}

	resp := dto.Response{
		Version: version,
		Data:    dto.Data(struct{ Data interface{} }{Data: message}),
	}

	jsonResp, _ := json.Marshal(resp)
	_, _ = w.Write(jsonResp)
	return w
}

// parseTelegramRequest handles incoming update from the Telegram web hook
func parseTelegramRequest(r *http.Request) (*telegramApi.Update, error) {
	var update telegramApi.Update

	if r.Method != http.MethodPost {
		log.Printf("unsupported method %s", r.Method)
		return nil, errors.New("unsupported method " + r.Method)
	}

	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		log.Printf("could not decode incoming update %s", err.Error())
		return nil, err
	}

	if update.UpdateID == 0 {
		log.Printf("invalid update id, got update id = 0")
		return nil, errors.New("invalid update id of 0 indicates failure to parse incoming update")
	}

	return &update, nil
}
