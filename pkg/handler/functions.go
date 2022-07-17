package handler

import (
	"encoding/json"
	"errors"
	telegramApi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"net/http"
	"os"
)

const telegramTokenEnv string = "TELEGRAM_BOT_TOKEN"

// HandleTelegramWebHook sends a message back to the chat with a punchline starting by the message provided by the user.
func HandleTelegramWebHook(w http.ResponseWriter, r *http.Request) {
	var bot, errBot = telegramApi.NewBotAPI(os.Getenv(telegramTokenEnv))
	if errBot != nil {
		responseMessage("Could not start bot", http.StatusInternalServerError, w)
		log.Fatalf("error could not start bot, %s", errBot)
		return
	}

	// Parse incoming request
	var update, err = parseTelegramRequest(r)

	if err != nil {
		responseMessage("Could not read update request", http.StatusUnprocessableEntity, w)
		log.Fatalf("error parsing update, %s", err.Error())
		return
	}

	if update.Message != nil {
		ActionHandler(update, bot)
	}

	responseMessage("Operation performed successfully", http.StatusOK, w)
	return
}

// parseTelegramRequest handles incoming update from the Telegram web hook
func parseTelegramRequest(r *http.Request) (*telegramApi.Update, error) {
	var update telegramApi.Update
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

// ResponseMessage responsible for returning a message about the status of the request
func responseMessage(m string, s int, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(s)

	resp := make(map[string]string)
	resp["message"] = m
	resp["version"] = os.Getenv("APP_VERSION")

	jsonResp, _ := json.Marshal(resp)
	_, _ = w.Write(jsonResp)
	return
}
