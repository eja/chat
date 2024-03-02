// Copyright (C) 2023-2024 by Ubaldo Porcheddu <ubaldo@eja.it>

package openai

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

	"github.com/eja/chat/internal/sys"
)

const asrModel = "whisper-1"

type typeASRResponse struct {
	Text string `json:"text"`
}

func ASR(filePath string, languageCode string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	writer.WriteField("model", asrModel)
	writer.WriteField("language", languageCode)
	filePart, err := writer.CreatePart(map[string][]string{
		"Content-Disposition": {"form-data; name=\"file\"; filename=\"" + filepath.Base(filePath) + ".ogg\""},
		"Content-Type":        {"audio/ogg"},
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

	request, err := http.NewRequest(http.MethodPost, "https://api.openai.com/v1/audio/transcriptions", &body)
	if err != nil {
		return "", err
	}
	request.Header.Set("Authorization", "Bearer "+sys.Options.OpenaiToken)
	request.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	var openAIResponse typeASRResponse
	err = json.Unmarshal(responseBody, &openAIResponse)
	if err != nil {
		return "", err
	}

	return openAIResponse.Text, nil
}
