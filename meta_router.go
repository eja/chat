// Copyright (C) 2023-2024 by Ubaldo Porcheddu <ubaldo@eja.it>

package main

import (
	"encoding/json"
	"net/http"
)

type typeMetaMessage struct {
	Entry []struct {
		Changes []struct {
			Value struct {
				Messages []struct {
					From string `json:"from"`
					Text *struct {
						Body string `json:"body"`
					} `json:"text,omitempty"`
					Audio *struct {
						ID string `json:"id"`
					} `json:"audio,omitempty"`
					ID string `json:"id"`
				} `json:"messages"`
			} `json:"value"`
		} `json:"changes"`
	} `json:"entry"`
}

func metaRouter(w http.ResponseWriter, r *http.Request) {
	if err := dbOpen(); err != nil {
		return
	}

	if r.Method == http.MethodGet {
		hubMode := r.URL.Query().Get("hub.mode")
		verifyToken := r.URL.Query().Get("hub.verify_token")
		if hubMode == "subscribe" && verifyToken == Options.MetaToken {
			w.Write([]byte(r.URL.Query().Get("hub.challenge")))
		} else {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Token verification error"))
		}
	} else {
		const platform = "meta"
		var metaMessage typeMetaMessage

		err := json.NewDecoder(r.Body).Decode(&metaMessage)
		if err != nil {
			errMessage := "Error decoding request body"
			http.Error(w, errMessage, http.StatusBadRequest)
			logWarn(errMessage)
			return
		}

		logTrace("FB incoming message", metaMessage)
		if len(metaMessage.Entry) > 0 {
			changes := metaMessage.Entry[0].Changes
			if len(changes) > 0 {
				value := changes[0].Value

				if len(value.Messages) > 0 {
					message := value.Messages[0]
					userId := message.From
					chatId := message.ID

					user, err := dbUserGet(userId)
					if err == nil && user != nil {
						if err := metaSendStatus(chatId, "read"); err != nil {
							logWarn("FB status", userId, chatId, err)
						}

						if dbNumber(user["welcome"]) < 1 {
							metaSendText(userId, translate(user["language"], "welcome"))
							dbUserUpdate(userId, "welcome", "1")
						}

						if message.Text != nil && message.Text.Body != "" {
							response, err := processText(userId, user["language"], message.Text.Body)
							if err != nil {
								response = translate(user["language"], "error")
								logWarn("FB", userId, chatId, err)
							}
							if err := metaSendText(userId, response); err != nil {
								logWarn("FB", userId, err)
							}
						} else if message.Audio != nil {
							if dbNumber(user["audio"]) > 0 {
								_, err := processAudio(
									platform,
									userId,
									user["language"],
									chatId,
									message.Audio.ID,
									dbNumber(user["audio"]) > 1,
								)
								if err != nil {
									logWarn("FB", userId, chatId, err)
									if err := metaSendText(userId, translate(user["language"], "error")); err != nil {
										logWarn("FB", userId, chatId, err)
									}
								}
							}
						} else {
							if err := metaSendText(userId, translate(user["language"], "audio_disabled")); err != nil {
								logWarn("FB", userId, chatId, err)
							}
						}
					} else {
						if err := metaSendText(userId, translate("", "user_unknown")); err != nil {
							logWarn("FB", userId, chatId, err)
						}
					}
				}
			}
		}
		w.WriteHeader(http.StatusOK)
	}
}
