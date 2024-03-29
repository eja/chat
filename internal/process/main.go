// Copyright (C) 2023-2024 by Ubaldo Porcheddu <ubaldo@eja.it>

package process

import (
	"fmt"
	"regexp"

	"github.com/eja/chat/internal/db"
	"github.com/eja/chat/internal/ff"
	"github.com/eja/chat/internal/google"
	"github.com/eja/chat/internal/i18n"
	"github.com/eja/chat/internal/meta"
	"github.com/eja/chat/internal/openai"
	"github.com/eja/chat/internal/sys"
	"github.com/eja/chat/internal/telegram"
)

const maxAudioInputTime = 60
const maxAudioOutputSize = 50 * 1000

func FilterLanguage(response string, language string) (string, string) {
	matchPattern := regexp.MustCompile(`\[\w{2}\]\s*$`)
	languagePattern := regexp.MustCompile(`\[(.*?)\]`)

	response = matchPattern.ReplaceAllStringFunc(response, func(code string) string {
		language = languagePattern.ReplaceAllString(code, "$1")
		return ""
	})

	return response, language
}

func Text(userId string, language string, text string) (string, error) {
	response, err := Chat(userId, text, language)
	if err == nil {
		response, _ = FilterLanguage(response, language)
	}
	return response, err
}

func Audio(platform string, userId string, language string, chatId string, mediaId string, tts bool) (string, error) {
	var response string
	var transcript string
	var err error

	mediaPath := fmt.Sprintf("%s/%s", sys.Options.MediaPath, mediaId)

	fileAudioInput := mediaPath + ".original.audio.in"
	if platform == "meta" {
		if err := meta.MediaGet(mediaId, fileAudioInput); err != nil {
			return "", err
		}
	}
	if platform == "telegram" {
		if err := telegram.MediaGet(mediaId, fileAudioInput); err != nil {
			return "", err
		}
	}
	if platform == "pbx" {
		fileAudioInput = mediaPath
	}

	if sys.Options.GoogleCredentials != "" {
		fileGoogleInput := mediaPath + ".google.audio.in"
		probeInput, err := ff.ProbeAudio(fileAudioInput)
		if err != nil {
			return "", err
		}
		duration := db.Number(probeInput["duration"])
		if duration > maxAudioInputTime {
			return i18n.Translate(language, "audio_input_limit"), nil
		}
		if err := ff.MpegAudioGoogle(fileAudioInput, fileGoogleInput); err != nil {
			return "", err
		}
		transcript, err = google.ASR(fileGoogleInput, i18n.LanguageCodeToLocale(language))
		if err != nil {
			return "", err
		}
	} else {
		transcript, err = openai.ASR(fileAudioInput, language)
		if err != nil {
			return "", err
		}
	}

	response, err = Chat(userId, transcript, language)
	if err != nil {
		return "", err
	}

	if !tts || len(response) > maxAudioOutputSize {
		return response, nil
	}

	response, ttsLanguage := FilterLanguage(response, language)

	fileAudioOutput := mediaPath + ".audio.out"

	if sys.Options.GoogleCredentials != "" {
		if err = google.TTS(fileAudioOutput, response, i18n.LanguageCodeToLocale(ttsLanguage)); err != nil {
			return "", err
		}
	} else {
		if err = openai.TTS(fileAudioOutput, response); err != nil {
			return "", err
		}
	}

	if platform == "meta" {
		if err := meta.SendAudio(userId, fileAudioOutput); err != nil {
			return "", err
		}
		response = ""
	}
	if platform == "telegram" {
		if err := telegram.SendAudio(chatId, fileAudioOutput, response); err != nil {
			return "", err
		}
		response = ""
	}
	if platform == "pbx" {
		response = fileAudioOutput
	}

	return response, nil
}
