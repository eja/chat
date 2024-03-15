// Copyright (C) 2023-2024 by Ubaldo Porcheddu <ubaldo@eja.it>

package sys

import (
	"flag"

	"github.com/eja/tibula/sys"
)

const Name = "chat"
const Version = "2.3.7"

var Options typeConfigChat

type typeConfigChat struct {
	sys.TypeConfig
	MediaPath         string `json:"media_path,omitempty"`
	GoogleCredentials string `json:"google_credentials,omitempty"`
	MetaUrl           string `json:"meta_url,omitempty"`
	MetaUser          string `json:"meta_user,omitempty"`
	MetaAuth          string `json:"meta_auth,omitempty"`
	MetaToken         string `json:"meta_token,omitempty"`
	TelegramToken     string `json:"telegram_token,omitempty"`
	OpenaiToken       string `json:"openai_token,omitempty"`
	PbxToken          string `json:"pbx_token,omitempty"`
}

func Configure() error {
	flag.StringVar(&Options.MediaPath, "media-path", "/tmp/", "Media temporary folder")
	flag.StringVar(&Options.GoogleCredentials, "google-credentials", "", "Google application credentials file path")
	flag.StringVar(&Options.MetaUrl, "meta-url", "", "Meta graph api url")
	flag.StringVar(&Options.MetaUser, "meta-user", "", "Meta user id")
	flag.StringVar(&Options.MetaAuth, "meta-auth", "", "Meta auth")
	flag.StringVar(&Options.MetaToken, "meta-token", "", "Meta token")
	flag.StringVar(&Options.TelegramToken, "telegram-token", "", "Telegram token")
	flag.StringVar(&Options.OpenaiToken, "openai-token", "", "OpenAI token")
	flag.StringVar(&Options.PbxToken, "pbx-token", "", "PBX token")

	if err := sys.Configure(); err != nil {
		return err
	}
	Options.TypeConfig = sys.Options

	if sys.Commands.Start && sys.Options.ConfigFile == "" {
		sys.Options.ConfigFile = Name + ".json"
	}

	if sys.Options.ConfigFile != "" {
		if err := sys.ConfigRead(sys.Options.ConfigFile, &Options); err != nil {
			return err
		}
	}

	return nil
}
