// Copyright (C) 2023-2024 by Ubaldo Porcheddu <ubaldo@eja.it>

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"
)

func metaRequest(method string, url string, body interface{}, headers map[string]string) ([]byte, error) {
	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			return nil, fmt.Errorf("encoding request body: %w", err)
		}
	}

	req, err := http.NewRequest(method, url, &buf)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", Options.MetaAuth))
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("request failed with status: %d", resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}

	return data, nil
}

func metaPost(data interface{}) error {
	url := fmt.Sprintf("%s/%s/messages", Options.MetaUrl, Options.MetaUser)
	headers := map[string]string{
		"Content-Type": "application/json",
	}

	_, err := metaRequest("POST", url, data, headers)
	return err
}

func metaMediaGet(mediaId string, fileName string) error {
	url := fmt.Sprintf("%s/%s/", Options.MetaUrl, mediaId)
	headers := map[string]string{}

	responseData, err := metaRequest("GET", url, nil, headers)
	if err != nil {
		return err
	}
	var data struct {
		URL string `json:"url"`
	}
	if err := json.Unmarshal(responseData, &data); err != nil {
		return fmt.Errorf("decoding response: %w", err)
	}

	responseData, err = metaRequest("GET", data.URL, nil, headers)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(fileName, responseData, 0644); err != nil {
		return fmt.Errorf("writing file: %w", err)
	}

	log.Printf("Media content saved to: %s", fileName)
	return nil
}

func metaMediaUpload(fileName string, fileType string) (mediaId string, err error) {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return "", fmt.Errorf("reading file: %w", err)
	}

	var b bytes.Buffer
	mw := multipart.NewWriter(&b)
	part, err := mw.CreateFormFile("file", filepath.Base(fileName))
	if err != nil {
		return "", fmt.Errorf("creating form part: %w", err)
	}
	part.Write(data)
	mw.WriteField("type", fileType)
	mw.WriteField("messaging_product", "whatsapp")
	mw.Close()

	responseData, err := metaRequest(
		"POST",
		fmt.Sprintf("%s/%s/media/", Options.MetaUrl, Options.MetaUser),
		b,
		map[string]string{
			"Content-Type": mw.FormDataContentType(),
		},
	)
	if err != nil {
		return "", err
	}

	var response struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(responseData, &response); err != nil {
		return "", fmt.Errorf("parsing response: %w", err)
	}

	log.Printf("meta media upload %s %s\n", fileName, fileType)
	return response.ID, nil
}

func metaSendText(phone string, text string) error {
	messageData := map[string]interface{}{
		"messaging_product": "whatsapp",
		"preview_url":       false,
		"recipient_type":    "individual",
		"to":                phone,
		"type":              "text",
		"text": map[string]interface{}{
			"body": text,
		},
	}

	_, err := metaRequest(
		"POST",
		fmt.Sprintf("%s/%s/messages", Options.MetaUrl, Options.MetaUser),
		messageData,
		nil,
	)
	return err
}

func metaSendStatus(messageId string, status string) error {
	statusData := map[string]interface{}{
		"messaging_product": "whatsapp",
		"message_id":        messageId,
		"status":            status,
	}

	_, err := metaRequest(
		"POST",
		fmt.Sprintf("%s/%s/messages/%s/status", Options.MetaUrl, Options.MetaUser, messageId),
		statusData,
		nil,
	)
	return err
}

func metaReaction(recipient string, messageId string, emoji string) (string, error) {
	reactionData := map[string]interface{}{
		"messaging_product": "whatsapp",
		"recipient_type":    "individual",
		"to":                recipient,
		"type":              "reaction",
		"reaction": map[string]interface{}{
			"message_id": messageId,
			"emoji":      emoji,
		},
	}

	responseData, err := metaRequest(
		"POST",
		fmt.Sprintf("%s/%s/messages", Options.MetaUrl, Options.MetaUser),
		reactionData,
		nil,
	)
	if err != nil {
		return "", err
	}

	return string(responseData), nil
}

func metaSendAudio(phone string, mediaFile string) error {
	mediaPath := filepath.Join(Options.MediaPath, phone)
	fileAudioOutput := mediaPath + ".audio.meta.out"

	probeOutput, err := ffprobeAudio(mediaFile)
	if err != nil {
		return fmt.Errorf("probing audio: %w", err)
	}
	if probeOutput["codecName"] == "opus" && probeOutput["sampleRate"] == "48000" && probeOutput["channelLayout"] == "mono" {
		fileAudioOutput = mediaFile
	} else {
		err = ffmpegAudioMeta(mediaFile, fileAudioOutput)
		if err != nil {
			return fmt.Errorf("converting audio: %w", err)
		}
	}

	mediaUploadId, err := metaMediaUpload(fileAudioOutput, "audio/ogg")
	if err != nil {
		return fmt.Errorf("uploading audio: %w", err)
	}

	messageData := map[string]interface{}{
		"messaging_product": "whatsapp",
		"recipient_type":    "individual",
		"to":                phone,
		"type":              "audio",
		"audio": map[string]interface{}{
			"id": mediaUploadId,
		},
	}
	_, err = metaRequest("POST", fmt.Sprintf("%s/%s/messages", Options.MetaUrl, Options.MetaUser), messageData, nil)
	return err
}
