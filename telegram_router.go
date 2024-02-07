// Copyright (C) 2023-2024 by Ubaldo Porcheddu <ubaldo@eja.it>

package main

import (
	"encoding/json"
	"fmt"
	"github.com/eja/tibula/db"
	"net/http"
)

type typeTelegramMessage struct {
	Message struct {
		From struct {
			Id           int    `json:"id"`
			LanguageCode string `json:"language_code"`
		} `json:"from"`
		Chat struct {
			Id int `json:"id"`
		} `json:"chat"`
		Text  string `json:"text,omitempty"`
		Voice struct {
			FileId string `json:"file_id"`
		} `json:"voice,omitempty"`
		Context struct {
			Forwarded bool `json:"forwarded"`
		} `json:"context,omitempty"`
	} `json:"message"`
}

func telegramRouter(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		const platform = "telegram"
		var telegramMessage typeTelegramMessage

		err := json.NewDecoder(r.Body).Decode(&telegramMessage)
		if err != nil {
			errMessage := "Error decoding request body"
			http.Error(w, errMessage, http.StatusBadRequest)
			logWarn(errMessage)
			return
		}

		logTrace("TG incoming message", telegramMessage)
		userId := fmt.Sprintf("TG.%d", telegramMessage.Message.From.Id)
		chatId := fmt.Sprintf("%d", telegramMessage.Message.Chat.Id)
		chatLanguage := telegramMessage.Message.From.LanguageCode

		user, _ := dbUserGet(userId)
		if db.Number(user["ejaId"]) > 0 {
			if db.Number(user["welcome"]) < 1 {
				telegramSendText(
					chatId,
					translate(chatLanguage, "welcome"),
				)
				dbUserUpdate(userId, "welcome", "1")
			}

			if text := telegramMessage.Message.Text; text != "" {
				response, err := processText(userId, user["language"], text)
				if err != nil {
					response = translate(user["language"], "process_text_error")
					logWarn("TG", userId, chatId, err)
				}
				if err := telegramSendText(
					chatId,
					response,
				); err != nil {
					logWarn("TG", userId, chatId, err)
				}
			}

			if voice := telegramMessage.Message.Voice; voice.FileId != "" {
				if db.Number(user["audio"]) > 0 {
					_, err := processAudio(
						platform,
						userId,
						user["language"],
						chatId,
						voice.FileId,
						db.Number(user["audio"]) > 1,
					)
					if err != nil {
						logWarn("TG", userId, chatId, err)
						if err := telegramSendText(chatId, translate(chatLanguage, "process_audio_error")); err != nil {
							logWarn("TG", userId, chatId, err)
						}
					}
				} else {
					if err := telegramSendText(
						chatId,
						translate(user["language"], "audio_disabled"),
					); err != nil {
						logWarn("TG", userId, chatId, err)
					}
				}
			}
		} else {
			if err := telegramSendText(chatId, translate(chatLanguage, "user_unknown")); err != nil {
				logWarn("TG", userId, chatId, err)
			}
		}
	}
}
