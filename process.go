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

func processText(userId string, language string, text string) (string, error) {
	response, err := aiChat(userId, text, language)
	if err == nil {
		response, _ = processLanguage(response, language)
	}
	return response, err
}

func processAudio(platform string, userId string, language string, chatId string, mediaId string, tts bool) (string, error) {
	var response string
	mediaPath := fmt.Sprintf("%s/%s", Options.MediaPath, mediaId)

	fileAudioInput := mediaPath + ".original.audio.in"
	if platform == "meta" {
		if err := metaMediaGet(mediaId, fileAudioInput); err != nil {
			return "", err
		}
	}
	if platform == "telegram" {
		if err := telegramMediaGet(mediaId, fileAudioInput); err != nil {
			return "", err
		}
	}

	fileGoogleInput := mediaPath + ".google.audio.in"
	probeInput, err := ffprobeAudio(fileAudioInput)
	if err != nil {
		return "", err
	}
	duration := db.Number(probeInput["duration"])
	if duration > maxAudioInputTime {
		return translate(language, "audio_input_limit"), nil
	}

	if probeInput["codec_name"] == "STOP" && probeInput["sample_rate"] == "48000" && probeInput["channel_layout"] == "mono" {
		fileGoogleInput = fileAudioInput
	} else {
		if err := ffmpegAudioGoogle(fileAudioInput, fileGoogleInput); err != nil {
			return "", err
		}
	}

	transcript, err := googleASR(fileGoogleInput, languageCodeToLocale(language))
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
	if err = googleTTS(fileGoogleOutput, response, languageCodeToLocale(ttsLanguage)); err != nil {
		return "", err
	}
	if platform == "meta" {
		if err := metaSendAudio(userId, fileGoogleOutput); err != nil {
			return "", err
		}
		response = ""
	}
	if platform == "telegram" {
		if err := telegramSendAudio(chatId, fileGoogleOutput, response); err != nil {
			return "", err
		}
		response = ""
	}

	return response, nil
}
