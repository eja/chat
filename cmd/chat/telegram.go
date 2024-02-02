// Copyright (C) 2023-2024 by Ubaldo Porcheddu <ubaldo@eja.it>

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

func TelegramRouter(w http.ResponseWriter, r *http.Request) {
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

	}
}

func TelegramSendText(chatID string, text string) error {
	sendMessageURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", chatOptions.TelegramToken)

	payload := map[string]string{
		"chat_id": chatID,
		"text":    text,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error encoding payload: %v", err)
	}

	resp, err := http.Post(sendMessageURL, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return fmt.Errorf("error sending text message: %v", err)
	}
	defer resp.Body.Close()

	return nil
}

func TelegramMediaGet(fileID string, fileName string) error {
	getFileURL := fmt.Sprintf("https://api.telegram.org/bot%s/getFile?file_id=%s", chatOptions.TelegramToken, fileID)

	resp, err := http.Get(getFileURL)
	if err != nil {
		return fmt.Errorf("error retrieving file path: %v", err)
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return fmt.Errorf("error parsing JSON response: %v", err)
	}

	if ok, result, filePath := data["ok"].(bool), data["result"].(map[string]interface{}), data["result"].(map[string]interface{})["file_path"].(string); ok && filePath != "" {
		filePath := result["file_path"].(string)
		fileURL := fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", chatOptions.TelegramToken, filePath)

		fileResponse, err := http.Get(fileURL)
		if err != nil {
			return fmt.Errorf("error downloading file: %v", err)
		}
		defer fileResponse.Body.Close()

		file, err := os.Create(fileName)
		if err != nil {
			return fmt.Errorf("error creating file: %v", err)
		}
		defer file.Close()

		_, err = io.Copy(file, fileResponse.Body)
		if err != nil {
			return fmt.Errorf("error writing to file: %v", err)
		}
	} else {
		return fmt.Errorf("failed to retrieve file path")
	}

	return nil
}

func TelegramSendAudio(chatID string, fileName string, caption string) error {
	sendVoiceURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendVoice", chatOptions.TelegramToken)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	writer.WriteField("chat_id", chatID)

	if caption != "" {
		writer.WriteField("caption", caption)
	}

	file, err := os.Open(fileName)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	part, err := writer.CreateFormFile("voice", fileName)
	if err != nil {
		return fmt.Errorf("error creating form file: %v", err)
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return fmt.Errorf("error copying file to form file: %v", err)
	}

	writer.Close()

	resp, err := http.Post(sendVoiceURL, writer.FormDataContentType(), body)
	if err != nil {
		return fmt.Errorf("error sending audio: %v", err)
	}
	defer resp.Body.Close()

	return nil
}
