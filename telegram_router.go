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
			errMessage := "Error decoding request body"
			http.Error(w, errMessage, http.StatusBadRequest)
			Debug(errMessage)
			return
		}

		Trace("TG incoming message", telegramMessage)
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
				response, err := ProcessText(userId, user["language"], text)
				if err != nil {
					response = Translate(user["language"], "process_text_error")
					Error("TG", userId, chatId, err)
				}
				if err := TelegramSendText(
					chatId,
					response,
				); err != nil {
					Error("TG", userId, chatId, err)
				}
			}

			if voice := telegramMessage.Message.Voice; voice.FileId != "" {
				if db.Number(user["audio"]) > 0 {
					_, err := ProcessAudio(
						platform,
						userId,
						user["language"],
						chatId,
						voice.FileId,
						db.Number(user["audio"]) > 1,
					)
					if err != nil {
						Error("TG", userId, chatId, err)
						if err := TelegramSendText(chatId, Translate(chatLanguage, "process_audio_error")); err != nil {
							Error("TG", userId, chatId, err)
						}
					}
				} else {
					if err := TelegramSendText(
						chatId,
						Translate(user["language"], "audio_disabled"),
					); err != nil {
						Error("TG", userId, chatId, err)
					}
				}
			}
		} else {
			if err := TelegramSendText(chatId, Translate(chatLanguage, "user_unknown")); err != nil {
				Error("TG", userId, chatId, err)
			}
		}
	}
}
