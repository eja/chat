// Copyright (C) 2023-2024 by Ubaldo Porcheddu <ubaldo@eja.it>

package main

import (
	"github.com/eja/tibula/sys"
)

func chatWizard() error {

	chatOptions.GoogleCredentials = sys.WizardPrompt("Google Application Credentials file path")
	chatOptions.MetaUrl = sys.WizardPrompt("Meta graph api url")
	chatOptions.MetaUser = sys.WizardPrompt("Meta user id")
	chatOptions.MetaAuth = sys.WizardPrompt("Meta auth")
	chatOptions.MetaToken = sys.WizardPrompt("Meta token")
	chatOptions.TelegramToken = sys.WizardPrompt("Telegram token")

	configFile := chatOptions.ConfigFile
	chatOptions.ConfigFile = ""
	return sys.ConfigWrite(configFile, &chatOptions)
}
