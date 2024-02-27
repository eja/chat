// Copyright (C) 2023-2024 by Ubaldo Porcheddu <ubaldo@eja.it>

package process

import (
	"fmt"
	"github.com/eja/chat/internal/db"
	"github.com/eja/chat/internal/i18n"
	"github.com/eja/chat/internal/log"
	"github.com/eja/chat/internal/openai"
	"regexp"
	"strings"
	"time"
)

const historyTimeout = 300
const model = "gpt-3.5-turbo-0125"

var history map[string][]typeAiMessage
var historyTime map[string]time.Time
var historyInit bool

type typeAiMessage = openai.TypeMessage

func Chat(userId, message, language string) (string, error) {
	log.Trace(userId, message, language)
	if !historyInit {
		history = make(map[string][]typeAiMessage)
		historyTime = make(map[string]time.Time)
		historyInit = true
	}
	var response, system string

	if rows, err := db.SystemPrompt(); err != nil {
		return "", err
	} else {
		for _, row := range rows {
			system += row["prompt"] + "\n"
		}
	}
	system += fmt.Sprintf("The user usually speaks in %s, so please answer in that language or the language of the question if not instructed otherwise.\n", i18n.LanguageCodeToInternal(language)) +
		"Always append a new line containing only the language code between square brackets that you have used to answer the question at the end of your response, like this: \n[en]\n" +
		""

	if strings.HasPrefix(message, "/") {
		if message == "/help" {
			response = i18n.Translate(language, "help")
		}

		if message == "/reset" {
			delete(history, userId)
			response = i18n.Translate(language, "reset")
		}

		if message == "/audio on" {
			user, err := db.UserGet(userId)
			if err != nil {
				return "", err
			}
			if db.Number(user["audio"]) > 0 {
				err := db.UserUpdate(userId, "audio", "2")
				if err != nil {
					return "", err
				}
				response = i18n.Translate(language, "audio_on")
			} else {
				response = i18n.Translate(language, "audio_disabled")
			}
		}

		if message == "/audio off" {
			user, err := db.UserGet(userId)
			if err != nil {
				return "", err
			}
			if db.Number(user["audio"]) > 0 {
				err := db.UserUpdate(userId, "audio", "1")
				if err != nil {
					return "", err
				}
				response = i18n.Translate(language, "audio_off")
			} else {
				response = i18n.Translate(language, "audio_disabled")
			}
		}

		if matched, _ := regexp.MatchString(`^/[a-zA-Z]{2}$`, message); matched {
			language = message[1:]
			err := db.UserUpdate(userId, "language", language)
			if err != nil {
				return "", err
			}
			response = i18n.Translate(language, "welcome")
			delete(history, userId)
		}
	}

	if response == "" {
		if hist, ok := history[userId]; ok && len(hist) > 0 && (time.Now().Sub(historyTime[userId]).Seconds() < historyTimeout) {
			history[userId] = append(history[userId], typeAiMessage{
				Role:    "user",
				Content: message,
			})
		} else {
			history[userId] = []typeAiMessage{
				{Role: "system", Content: system},
				{Role: "user", Content: message},
			}
		}

		assistant, err := openai.Request(model, history[userId])
		if err != nil {
			return "", err
		}
		historyTime[userId] = time.Now()
		history[userId] = append(history[userId], typeAiMessage{
			Role:    "assistant",
			Content: assistant,
		})
		response = assistant
	}

	return response, nil
}
