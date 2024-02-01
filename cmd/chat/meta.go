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

func MetaRouter(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		hubMode := r.URL.Query().Get("hub.mode")
		verifyToken := r.URL.Query().Get("hub.verify_token")
		if hubMode == "subscribe" && verifyToken == chatOptions.MetaToken {
			w.Write([]byte(r.URL.Query().Get("hub.challenge")))
		} else {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Token verification error"))
		}
	} else {
		//
	}
}

func MetaRequest(method string, url string, body interface{}, headers map[string]string) ([]byte, error) {
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

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", chatOptions.MetaAuth))
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

func MetaPost(data interface{}) error {
	url := fmt.Sprintf("%s/%s/messages", chatOptions.MetaUrl, chatOptions.MetaUser)
	headers := map[string]string{
		"Content-Type": "application/json",
	}

	_, err := MetaRequest("POST", url, data, headers)
	return err
}

func MetaMediaGet(mediaId string, fileName string) error {
	url := fmt.Sprintf("%s/%s/", chatOptions.MetaUrl, mediaId)
	headers := map[string]string{}

	responseData, err := MetaRequest("GET", url, nil, headers)
	if err != nil {
		return err
	}
	var data struct {
		URL string `json:"url"`
	}
	if err := json.Unmarshal(responseData, &data); err != nil {
		return fmt.Errorf("decoding response: %w", err)
	}

	responseData, err = MetaRequest("GET", data.URL, nil, headers)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(fileName, responseData, 0644); err != nil {
		return fmt.Errorf("writing file: %w", err)
	}

	log.Printf("Media content saved to: %s", fileName)
	return nil
}

func MetaMediaUpload(fileName string, fileType string) (mediaId string, err error) {
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

	responseData, err := MetaRequest(
		"POST",
		fmt.Sprintf("%s/%s/media/", chatOptions.MetaUrl, chatOptions.MetaUser),
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

func MetaSendText(phone string, text string) error {
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

	_, err := MetaRequest(
		"POST",
		fmt.Sprintf("%s/%s/messages", chatOptions.MetaUrl, chatOptions.MetaUser),
		messageData,
		nil,
	)
	return err
}

func MetaSendStatus(messageId string, status string) error {
	statusData := map[string]interface{}{
		"messaging_product": "whatsapp",
		"message_id":        messageId,
		"status":            status,
	}

	_, err := MetaRequest(
		"POST",
		fmt.Sprintf("%s/%s/messages/%s/status", chatOptions.MetaUrl, chatOptions.MetaUser, messageId),
		statusData,
		nil,
	)
	return err
}

func MetaReaction(recipient string, messageId string, emoji string) (string, error) {
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

	responseData, err := MetaRequest(
		"POST",
		fmt.Sprintf("%s/%s/messages", chatOptions.MetaUrl, chatOptions.MetaUser),
		reactionData,
		nil,
	)
	if err != nil {
		return "", err
	}

	return string(responseData), nil
}

func MetaSendAudio(phone string, mediaFile string) error {
	mediaPath := filepath.Join(chatOptions.MediaPath, phone)
	fileAudioOutput := mediaPath + ".audio.meta.out"

	probeOutput, err := FFProbeAudio(mediaFile) // Assuming ffProbeAudio returns a struct with codec_name, sample_rate, and channel_layout
	if err != nil {
		return fmt.Errorf("probing audio: %w", err)
	}
	if probeOutput["codecName"] == "opus" && probeOutput["sampleRate"] == "48000" && probeOutput["channelLayout"] == "mono" {
		fileAudioOutput = mediaFile
	} else {
		err = FFAudioMeta(mediaFile, fileAudioOutput) // Assuming ffAudioMeta converts audio format if needed
		if err != nil {
			return fmt.Errorf("converting audio: %w", err)
		}
	}

	mediaUploadId, err := MetaMediaUpload(fileAudioOutput, "audio/ogg")
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
	_, err = MetaRequest("POST", fmt.Sprintf("%s/%s/messages", chatOptions.MetaUrl, chatOptions.MetaUser), messageData, nil)
	return err
}
