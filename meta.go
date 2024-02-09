// Copyright (C) 2023-2024 by Ubaldo Porcheddu <ubaldo@eja.it>

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func metaRequest(method string, url string, body interface{}, contentType string) ([]byte, error) {
	var buf bytes.Buffer
	if contentType == "json" && body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			return nil, fmt.Errorf("encoding request body: %w", err)
		}
	}

	if strings.Contains(contentType, "multipart") && body != nil {
		buf = body.(bytes.Buffer)
	}

	req, err := http.NewRequest(method, url, &buf)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", Options.MetaAuth))
	if contentType == "json" {
		req.Header.Set("Content-Type", "application/json")
	}
	if strings.Contains(contentType, "multipart") {
		req.Header.Set("Content-Type", contentType)
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
	_, err := metaRequest("POST", url, data, "json")
	return err
}

func metaGet(url string) ([]byte, error) {
	return metaRequest("GET", url, nil, "")
}

func metaMediaGet(mediaId string, fileName string) error {
	url := fmt.Sprintf("%s/%s/", Options.MetaUrl, mediaId)

	responseData, err := metaGet(url)
	if err != nil {
		return err
	}
	var data struct {
		URL string `json:"url"`
	}
	if err := json.Unmarshal(responseData, &data); err != nil {
		return fmt.Errorf("decoding response: %w", err)
	}

	responseData, err = metaGet(data.URL)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(fileName, responseData, 0644); err != nil {
		return fmt.Errorf("writing file: %w", err)
	}

	logTrace("meta media content saved to: %s", fileName)
	return nil
}

func metaMediaUpload(fileName string, fileType string) (mediaId string, err error) {
	file, err := os.Open(fileName)
	if err != nil {
		return "", fmt.Errorf("reading file: %w", err)
	}
	defer file.Close()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	contentType := writer.FormDataContentType()
	writer.WriteField("type", fileType)
	writer.WriteField("messaging_product", "whatsapp")
	filePart, err := writer.CreatePart(map[string][]string{
		"Content-Disposition": {"form-data; name=\"file\"; filename=\"" + fileName + "\""},
		"Content-Type":        {fileType},
	})
	if err != nil {
		return "", fmt.Errorf("creating form file")
	}
	_, err = io.Copy(filePart, file)
	if err != nil {
		return "", fmt.Errorf("copying file into part")
	}
	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("closing form writer: %w", err)
	}

	responseData, err := metaRequest(
		"POST",
		fmt.Sprintf("%s/%s/media", Options.MetaUrl, Options.MetaUser),
		body,
		contentType,
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

	logTrace("meta media upload %s %s\n", fileName, fileType)
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

	return metaPost(messageData)
}

func metaSendStatus(messageId string, status string) error {
	statusData := map[string]interface{}{
		"messaging_product": "whatsapp",
		"message_id":        messageId,
		"status":            status,
	}

	return metaPost(statusData)
}

func metaReaction(recipient string, messageId string, emoji string) error {
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

	return metaPost(reactionData)
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
	return metaPost(messageData)
}
