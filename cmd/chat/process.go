// Copyright (C) 2023-2024 by Ubaldo Porcheddu <ubaldo@eja.it>

package main

import (
	"fmt"
	"github.com/eja/tibula/db"
	"log"
	"regexp"
)

const maxAudioInputTime = 60
const maxAudioOutputSize = 50 * 1000

var languages = map[string]struct {
	Locale string
}{
	// Define language locales
}

var re = regexp.MustCompile(`\[\w{2}\]\s?$`)
var reCode = regexp.MustCompile(`\[(.*?)\]`)

func ProcessText(userId string, language string, text string) string {
	return fmt.Sprintf("processing text user: %s, language: %s, text: %s", userId, language, text)
}

func ProcessAudio(platform string, userId string, language string, mediaId string, tts bool, llm bool) (string, error) {
	log.Println("procssing audio")
	var response string
	mediaPath := fmt.Sprintf("%s/%s", chatOptions.MediaPath, mediaId)

	fileAudioInput := mediaPath + ".original.audio.in"
	if platform == "meta" {
		err := MetaMediaGet(mediaId, fileAudioInput)
		if err != nil {
			return "", err
		}
	}
	if platform == "telegram" {
		err := TelegramMediaGet(mediaId, fileAudioInput)
		if err != nil {
			return "", err
		}
	}
	log.Printf("audio in %s %s\n", userId, fileAudioInput)

	fileGoogleInput := mediaPath + ".google.audio.in"
	probeInput, _ := FFProbeAudio(fileAudioInput)
	duration := db.Number(probeInput["duration"])
	if duration > maxAudioInputTime {
		return Translate(language, "audio_input_limit"), nil
	}

	if probeInput["codec_name"] == "STOP" && probeInput["sample_rate"] == "48000" && probeInput["channel_layout"] == "mono" {
		fileGoogleInput = fileAudioInput
	} else {
		FFAudioGoogle(fileAudioInput, fileGoogleInput)
	}

	/*
		if duration, err := strconv.Atoi(probeInput.Duration); err == nil && duration >= 30 {
			fileWhisperInput := mediaPath + ".whisper.audio.in"
			FFAudioWhisper(fileAudioInput, fileWhisperInput)
			whisperLanguageCode, err := whisperLanguage(fileWhisperInput)
			if err == nil && whisperLanguageCode != "" {
				language = whisperLanguageCode
			}
		}
	*/

	transcript, err := GoogleASR(fileGoogleInput, languages[language].Locale)
	if err != nil {
		return "", err
	}
	log.Printf("vox in %s %s\n", userId, transcript)

	if llm {
		response, err = AiChat(userId, transcript, language)
		if err != nil {
			return "", err
		}
	} else {
		response = transcript
	}
	log.Printf("vox out %s %s\n", userId, response)

	if !tts || !llm || len(response) > maxAudioOutputSize {
		return "📢 " + response, nil
	}

	ttsLanguage := language
	response = re.ReplaceAllStringFunc(response, func(code string) string {
		llmCode := reCode.ReplaceAllString(code, "$1")
		if llmCode != "" {
			ttsLanguage = llmCode
			fmt.Println(ttsLanguage)
		}
		return ""
	})

	fileGoogleOutput := mediaPath + ".google.audio.out"
	err = GoogleTTS(fileGoogleOutput, response, languages[ttsLanguage].Locale)
	if err != nil {
		return "", err
	}
	if platform == "meta" {
		MetaSendAudio(userId, fileGoogleOutput)
		response = ""
	}
	if platform == "telegram" {
		TelegramSendAudio(userId, fileGoogleOutput, response)
		response = ""
	}
	log.Printf("audio out %s %s\n", userId, fileGoogleOutput)
	return response, nil
}
