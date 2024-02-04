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
	MediaPath         string `json:"media_path,omitempty"`
	GoogleCredentials string `json:"google_credentials,omitempty"`
	MetaUrl           string `json:"meta_url,omitempty"`
	MetaUser          string `json:"meta_user,omitempty"`
	MetaAuth          string `json:"meta_auth,omitempty"`
	MetaToken         string `json:"meta_token,omitempty"`
	TelegramToken     string `json:"telegram_token,omitempty"`
	OpenaiToken       string `json:"openai_token,omitempty"`
}

func Configure() error {
	flag.StringVar(&chatOptions.MediaPath, "media-path", "/tmp/", "Media temporary folder")
	flag.StringVar(&chatOptions.GoogleCredentials, "google-credentials", "google.json", "Google application credentials file path")
	flag.StringVar(&chatOptions.MetaUrl, "meta-url", "", "Meta graph api url")
	flag.StringVar(&chatOptions.MetaUser, "meta-user", "", "Meta user id")
	flag.StringVar(&chatOptions.MetaAuth, "meta-auth", "", "Meta auth")
	flag.StringVar(&chatOptions.MetaToken, "meta-token", "", "Meta token")
	flag.StringVar(&chatOptions.TelegramToken, "telegram-token", "", "Telegram token")
	flag.StringVar(&chatOptions.OpenaiToken, "openai-token", "", "OpenAI token")

	if err := sys.Configure(); err != nil {
		return err
	}
	chatOptions.TypeConfig = sys.Options

	if sys.Commands.Start && sys.Options.ConfigFile == "" {
		sys.Options.ConfigFile = Name + ".json"
	}

	if sys.Options.ConfigFile != "" {
		if err := sys.ConfigRead(sys.Options.ConfigFile, &chatOptions); err != nil {
			return err
		}
	}

	return nil
}
