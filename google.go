// Copyright (C) 2023-2024 by Ubaldo Porcheddu <ubaldo@eja.it>

package main

import (
	"context"
	"io/ioutil"

	speech "cloud.google.com/go/speech/apiv1"
	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	"google.golang.org/api/option"
	speechpb "google.golang.org/genproto/googleapis/cloud/speech/v1"
	texttospeechpb "google.golang.org/genproto/googleapis/cloud/texttospeech/v1"
)

func GoogleASR(fileName string, language string) (string, error) {
	ctx := context.Background()

	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return "", err
	}

	client, err := speech.NewClient(ctx, option.WithCredentialsFile(Options.GoogleCredentials))
	if err != nil {
		return "", err
	}
	defer client.Close()

	config := &speechpb.RecognitionConfig{
		Encoding:        speechpb.RecognitionConfig_OGG_OPUS,
		SampleRateHertz: 48000,
		LanguageCode:    language,
	}

	audio := &speechpb.RecognitionAudio{
		AudioSource: &speechpb.RecognitionAudio_Content{
			Content: data, //base64.StdEncoding.EncodeToString(data),
		},
	}

	request := &speechpb.RecognizeRequest{
		Config: config,
		Audio:  audio,
	}

	resp, err := client.Recognize(ctx, request)
	if err != nil {
		return "", err
	}

	transcript := ""
	if len(resp.Results) > 0 && len(resp.Results[0].Alternatives) > 0 {
		transcript = resp.Results[0].Alternatives[0].Transcript
	}

	return transcript, nil
}

func GoogleTTS(fileName string, text string, language string) error {
	ctx := context.Background()

	client, err := texttospeech.NewClient(ctx, option.WithCredentialsFile(Options.GoogleCredentials))
	if err != nil {
		return err
	}
	defer client.Close()

	request := &texttospeechpb.SynthesizeSpeechRequest{
		Input: &texttospeechpb.SynthesisInput{
			InputSource: &texttospeechpb.SynthesisInput_Text{
				Text: text,
			},
		},
		Voice: &texttospeechpb.VoiceSelectionParams{
			LanguageCode: language,
			// SsmlGender:   texttospeechpb.SsmlVoiceGender_MALE,
		},
		AudioConfig: &texttospeechpb.AudioConfig{
			AudioEncoding: texttospeechpb.AudioEncoding_OGG_OPUS,
		},
	}

	resp, err := client.SynthesizeSpeech(ctx, request)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(fileName, resp.AudioContent, 0644)
	if err != nil {
		return err
	}

	return nil
}
