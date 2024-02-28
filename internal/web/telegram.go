// Copyright (C) 2023-2024 by Ubaldo Porcheddu <ubaldo@eja.it>

package web

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/eja/chat/internal/db"
	"github.com/eja/chat/internal/i18n"
	"github.com/eja/chat/internal/process"
	"github.com/eja/chat/internal/telegram"
	"github.com/eja/tibula/log"
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
	if err := db.Open(); err != nil {
		return
	}

	if r.Method == http.MethodPost {
		const platform = "telegram"
		var telegramMessage typeTelegramMessage

		err := json.NewDecoder(r.Body).Decode(&telegramMessage)
		if err != nil {
			errMessage := "Error decoding request body"
			http.Error(w, errMessage, http.StatusBadRequest)
			log.Warn("[TG]", errMessage)
			return
		}

		log.Trace("[TG]", "incoming message", telegramMessage)
		userId := fmt.Sprintf("TG.%d", telegramMessage.Message.From.Id)
		chatId := fmt.Sprintf("%d", telegramMessage.Message.Chat.Id)
		chatLanguage := telegramMessage.Message.From.LanguageCode

		user, err := db.UserGet(userId)
		if err == nil && user != nil {
			if db.Number(user["welcome"]) < 1 {
				telegram.SendText(
					chatId,
					i18n.Translate(chatLanguage, "welcome"),
				)
				db.UserUpdate(userId, "welcome", "1")
			}

			if text := telegramMessage.Message.Text; text != "" {
				response, err := process.Text(userId, user["language"], text)
				if err != nil {
					response = i18n.Translate(user["language"], "error")
					log.Warn("[TG]", userId, chatId, err)
				}
				if err := telegram.SendText(
					chatId,
					response,
				); err != nil {
					log.Warn("[TG]", userId, chatId, err)
				}
			}

			if voice := telegramMessage.Message.Voice; voice.FileId != "" {
				if db.Number(user["audio"]) > 0 {
					_, err := process.Audio(
						platform,
						userId,
						user["language"],
						chatId,
						voice.FileId,
						db.Number(user["audio"]) > 1,
					)
					if err != nil {
						log.Warn("[TG]", userId, chatId, err)
						if err := telegram.SendText(chatId, i18n.Translate(chatLanguage, "error")); err != nil {
							log.Warn("[TG]", userId, chatId, err)
						}
					}
				} else {
					if err := telegram.SendText(
						chatId,
						i18n.Translate(user["language"], "audio_disabled"),
					); err != nil {
						log.Warn("[TG]", userId, chatId, err)
					}
				}
			}
		} else {
			if err := telegram.SendText(chatId, i18n.Translate(chatLanguage, "user_unknown")); err != nil {
				log.Warn("[TG]", userId, chatId, err)
			}
		}
	}
}
