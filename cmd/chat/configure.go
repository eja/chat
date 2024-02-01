// Copyright (C) 2023-2024 by Ubaldo Porcheddu <ubaldo@eja.it>

package main

import (
	"flag"
	"github.com/eja/tibula/sys"
)

const Name = "chat"
const Version = "2.1.30"

var chatOptions TypeConfigChat

type TypeConfigChat struct {
	sys.TypeConfig
	GoogleCredentials string `json:"google_credentials,omitempty"`
	MetaUrl           string `json:"meta_url,omitempty"`
	MetaUser          string `json:"meta_user,omitempty"`
	MetaAuth          string `json:"meta_auth,omitempty"`
	MetaToken         string `json:"meta_token,omitempty"`
	TelegramToken     string `json:"telegram_token,omitempty"`
}

func Configure() error {
	flag.StringVar(&chatOptions.GoogleCredentials, "google-credentials", "google.json", "Google application credentials file path")
	flag.StringVar(&chatOptions.MetaUrl, "meta-url", "", "Meta graph api url")
	flag.StringVar(&chatOptions.MetaUser, "meta-user", "", "Meta user id")
	flag.StringVar(&chatOptions.MetaAuth, "meta-auth", "", "Meta auth")
	flag.StringVar(&chatOptions.MetaToken, "meta-token", "", "Meta token")
	flag.StringVar(&chatOptions.TelegramToken, "telegram-token", "", "Telegram token")

	if err := sys.Configure(); err != nil {
		return err
	}

	if chatOptions.ConfigFile != "" {
		if err := sys.ConfigRead(chatOptions.ConfigFile, &chatOptions); err != nil {
			return err
		}
	}

	return nil
}
