// Copyright (C) 2023-2024 by Ubaldo Porcheddu <ubaldo@eja.it>

package main

import (
	"github.com/eja/tibula/sys"
)

func chatWizard() error {

	chatOptions.MediaPath = sys.WizardPrompt("Media temporary folder path")
	chatOptions.GoogleCredentials = sys.WizardPrompt("Google Application Credentials file path")
	chatOptions.OpenaiToken = sys.WizardPrompt("OpenAI API key")
	chatOptions.TelegramToken = sys.WizardPrompt("Telegram token")
	chatOptions.MetaUrl = sys.WizardPrompt("Meta graph api url")
	if chatOptions.MetaUrl != "" {
		chatOptions.MetaUser = sys.WizardPrompt("Meta user id")
		chatOptions.MetaAuth = sys.WizardPrompt("Meta auth")
		chatOptions.MetaToken = sys.WizardPrompt("Meta token")
	}

	configFile := chatOptions.ConfigFile
	chatOptions.ConfigFile = ""
	return sys.ConfigWrite(configFile, &chatOptions)
}
