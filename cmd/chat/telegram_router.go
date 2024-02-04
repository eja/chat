// Copyright (C) 2023-2024 by Ubaldo Porcheddu <ubaldo@eja.it>

package main

import (
	"encoding/json"
	"fmt"
	"github.com/eja/tibula/db"
	"net/http"
)

type TypeTelegramMessage struct {
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

func TelegramRouter(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		const platform = "telegram"
		var telegramMessage TypeTelegramMessage

		err := json.NewDecoder(r.Body).Decode(&telegramMessage)
		if err != nil {
			http.Error(w, "Error decoding request body", http.StatusBadRequest)
			return
		}

		userId := fmt.Sprintf("TG.%d", telegramMessage.Message.From.Id)
		chatId := fmt.Sprintf("%d", telegramMessage.Message.Chat.Id)
		chatLanguage := telegramMessage.Message.From.LanguageCode

		user, _ := DbUserGet(userId)
		if db.Number(user["ejaId"]) > 0 {
			if db.Number(user["welcome"]) < 1 {
				TelegramSendText(
					chatId,
					Translate(chatLanguage, "welcome"),
				)
				DbUserUpdate(userId, "welcome", "1")
			}

			if text := telegramMessage.Message.Text; text != "" {
				TelegramSendText(
					fmt.Sprintf("%d", telegramMessage.Message.Chat.Id),
					ProcessText(userId, user["language"], text),
				)
			}

			if voice := telegramMessage.Message.Voice; voice.FileId != "" {
				llm := true
				if forwarded := telegramMessage.Message.Context.Forwarded; forwarded {
					llm = false
				}
				if llm {
					//?
				}
				if db.Number(user["audio"]) > 0 {
					response, err := ProcessAudio(
						platform,
						chatId,
						user["language"],
						voice.FileId,
						db.Number(user["audio"]) > 1,
						llm,
					)
					if err == nil && response != "" {
						TelegramSendText(chatId, response)
					}
				} else {
					TelegramSendText(
						chatId,
						Translate(user["language"], "audio_disabled"),
					)
				}
			}
		} else {
			TelegramSendText(chatId, Translate(chatLanguage, "user_unknown"))
		}
	}
}
