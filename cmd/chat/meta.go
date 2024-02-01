// Copyright (C) 2023-2024 by Ubaldo Porcheddu <ubaldo@eja.it>

package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
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

func MetaPost(data interface{}) error {
	url := fmt.Sprintf("%s/%s/messages", chatOptions.MetaUrl, chatOptions.MetaUser)

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+chatOptions.MetaAuth)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status)
	}

	return nil
}
