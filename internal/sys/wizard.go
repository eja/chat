// Copyright (C) 2023-2024 by Ubaldo Porcheddu <ubaldo@eja.it>

package sys

import (
	"embed"

	"github.com/eja/tibula/db"
	"github.com/eja/tibula/sys"
)

//go:embed all:assets
var chatDbAssets embed.FS

func Wizard() error {
	configFile := sys.Options.ConfigFile
	if err := sys.ConfigRead(configFile, &Options); err != nil {
		return err
	}

	Options.MediaPath = sys.WizardPrompt("Media temporary folder path")
	Options.GoogleCredentials = sys.WizardPrompt("Google Application Credentials file path")
	Options.OpenaiToken = sys.WizardPrompt("OpenAI API key")
	Options.TelegramToken = sys.WizardPrompt("Telegram token")
	Options.MetaUrl = sys.WizardPrompt("Meta graph api url")
	if Options.MetaUrl != "" {
		Options.MetaUser = sys.WizardPrompt("Meta user id")
		Options.MetaAuth = sys.WizardPrompt("Meta auth")
		Options.MetaToken = sys.WizardPrompt("Meta token")
	}

	db.Assets = chatDbAssets
	if err := db.Open(Options.DbType, Options.DbName, Options.DbUser, Options.DbPass, Options.DbHost, Options.DbPort); err != nil {
		return err
	}
	if err := db.Setup(""); err != nil {
		return err
	}

	Options.ConfigFile = ""
	return sys.ConfigWrite(configFile, &Options)
}
