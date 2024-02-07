// Copyright (C) 2023-2024 by Ubaldo Porcheddu <ubaldo@eja.it>

package main

import (
	"github.com/eja/tibula/sys"
)

func chatWizard() error {

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

	configFile := Options.ConfigFile
	Options.ConfigFile = ""
	return sys.ConfigWrite(configFile, &Options)
}
