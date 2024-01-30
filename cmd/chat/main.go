// Copyright (C) 2023-2024 by Ubaldo Porcheddu <ubaldo@eja.it>

package main

import (
	"flag"
	"fmt"
	"github.com/eja/tibula/sys"
	"github.com/eja/tibula/web"
	"log"
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

func main() {
	if err := Configure(); err != nil {
		log.Fatal(err)
	}
	chatOptions.TypeConfig = sys.Options
	if sys.Commands.Start {
		if chatOptions.DbName == "" && chatOptions.ConfigFile == "" {
			if err := sys.ConfigRead("config.json", &chatOptions); err != nil {
				log.Fatal("Config file missing or not enough parameters to continue.")
			}
		}
		if chatOptions.DbName == "" {
			log.Fatal("Database name/file is mandatory.")
		}
		if err := web.Start(); err != nil {
			log.Fatal("Cannot start the web service: ", err)
		}
	} else if sys.Commands.Wizard {
		if err := chatWizard(); err != nil {
			log.Fatal(err)
		}
	} else {
		Help()
	}
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

func Help() {
	fmt.Println("Copyright:", "2023-2024 by Ubaldo Porcheddu <ubaldo@eja.it>")
	fmt.Println("Version:", Version)
	fmt.Printf("Usage: %s [options]\n", Name)
	fmt.Println("\nOptions:")
	flag.PrintDefaults()
	fmt.Println()
}
