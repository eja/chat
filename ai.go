// Copyright (C) 2023-2024 by Ubaldo Porcheddu <ubaldo@eja.it>

package main

import (
	"fmt"
	"github.com/eja/tibula/db"
	"regexp"
	"strings"
	"time"
)

const historyTimeout = 300
const model = "gpt-3.5-turbo-0125"

var history map[string][]typeAiMessage
var historyTime map[string]time.Time
var historyInit bool

type typeAiMessage = typeOpenaiMessage

func aiChat(userId, message, language string) (string, error) {
	logTrace(userId, message, language)
	if !historyInit {
		history = make(map[string][]typeAiMessage)
		historyTime = make(map[string]time.Time)
		historyInit = true
	}
	var response, system string
	system = "You are a helpful assistant written and developed by eja.it.\n" +
		"Answer questions in a short but clear way.\n" +
		"eja.it is an open-source development company based between the wonderful islands of Sardinia and Tenerife.\n" +
		"Your name is Eja, and any reference to eja.it should point to you.\n" +
		"The engineer who wrote this bot is Ubaldo Porcheddu.\n" +
		"eja's most notable productions are:\n" +
		"eja.surf, an Android browser implementing all known privacy features\n" +
		"tibula, a Dynamic RDBMS with an Intuitive Web Interface and Flexible Customization\n" +
		"eja.tv, a portal to watch thousands of live IPTV channels from around the world\n" +
		"cronotopia, a portal to interact with historical events time and geo-referenced\n" +
		"When asked for information about you, eja, or eja.it, you must take into consideration the above data.\n" +
		fmt.Sprintf("The user usually speaks in %s, so please answer in that language or the language of the question if not instructed otherwise.\n", languageCodeToInternal(language)) +
		"Always append a new line containing only the language code between square brackets that you have used to answer the question at the end of your response, like this: \n[en]\n" +
		""

	if strings.HasPrefix(message, "/") {
		if message == "/help" {
			response = translate(language, "help")
		}

		if message == "/reset" {
			delete(history, userId)
			response = translate(language, "reset")
		}

		if message == "/audio on" {
			user, err := dbUserGet(userId)
			if err != nil {
				return "", err
			}
			if db.Number(user["audio"]) > 0 {
				err := dbUserUpdate(userId, "audio", "2")
				if err != nil {
					return "", err
				}
				response = translate(language, "audio_on")
			} else {
				response = translate(language, "audio_disabled")
			}
		}

		if message == "/audio off" {
			user, err := dbUserGet(userId)
			if err != nil {
				return "", err
			}
			if db.Number(user["audio"]) > 0 {
				err := dbUserUpdate(userId, "audio", "1")
				if err != nil {
					return "", err
				}
				response = translate(language, "audio_off")
			} else {
				response = translate(language, "audio_disabled")
			}
		}

		if matched, _ := regexp.MatchString(`^/[a-zA-Z]{2}$`, message); matched {
			language = message[1:]
			err := dbUserUpdate(userId, "language", language)
			if err != nil {
				return "", err
			}
			response = translate(language, "welcome")
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

		assistant, err := openaiRequest(model, history[userId])
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