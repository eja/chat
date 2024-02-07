// Copyright (C) 2023-2024 by Ubaldo Porcheddu <ubaldo@eja.it>

package main

import (
	"fmt"
	"github.com/eja/tibula/db"
	"regexp"
)

const maxAudioInputTime = 60
const maxAudioOutputSize = 50 * 1000

func processLanguage(response string, language string) (string, string) {
	matchPattern := regexp.MustCompile(`\[\w{2}\]\s*$`)
	languagePattern := regexp.MustCompile(`\[(.*?)\]`)

	response = matchPattern.ReplaceAllStringFunc(response, func(code string) string {
		language = languagePattern.ReplaceAllString(code, "$1")
		return ""
	})

	return response, language
}

func ProcessText(userId string, language string, text string) (string, error) {
	response, err := aiChat(userId, text, language)
	if err == nil {
		response, _ = processLanguage(response, language)
	}
	return response, err
}

func ProcessAudio(platform string, userId string, language string, chatId string, mediaId string, tts bool) (string, error) {
	var response string
	mediaPath := fmt.Sprintf("%s/%s", Options.MediaPath, mediaId)

	fileAudioInput := mediaPath + ".original.audio.in"
	if platform == "meta" {
		if err := MetaMediaGet(mediaId, fileAudioInput); err != nil {
			return "", err
		}
	}
	if platform == "telegram" {
		if err := TelegramMediaGet(mediaId, fileAudioInput); err != nil {
			return "", err
		}
	}

	fileGoogleInput := mediaPath + ".google.audio.in"
	probeInput, err := FFProbeAudio(fileAudioInput)
	if err != nil {
		return "", err
	}
	duration := db.Number(probeInput["duration"])
	if duration > maxAudioInputTime {
		return Translate(language, "audio_input_limit"), nil
	}

	if probeInput["codec_name"] == "STOP" && probeInput["sample_rate"] == "48000" && probeInput["channel_layout"] == "mono" {
		fileGoogleInput = fileAudioInput
	} else {
		if err := FFAudioGoogle(fileAudioInput, fileGoogleInput); err != nil {
			return "", err
		}
	}

	transcript, err := GoogleASR(fileGoogleInput, LanguageCodeToLocale(language))
	if err != nil {
		return "", err
	}

	response, err = aiChat(userId, transcript, language)
	if err != nil {
		return "", err
	}

	if !tts || len(response) > maxAudioOutputSize {
		return response, nil
	}

	response, ttsLanguage := processLanguage(response, language)

	fileGoogleOutput := mediaPath + ".google.audio.out"
	if err = GoogleTTS(fileGoogleOutput, response, LanguageCodeToLocale(ttsLanguage)); err != nil {
		return "", err
	}
	if platform == "meta" {
		if err := MetaSendAudio(userId, fileGoogleOutput); err != nil {
			return "", err
		}
		response = ""
	}
	if platform == "telegram" {
		if err := TelegramSendAudio(chatId, fileGoogleOutput, response); err != nil {
			return "", err
		}
		response = ""
	}

	return response, nil
}
