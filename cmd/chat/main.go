// Copyright (C) 2023-2024 by Ubaldo Porcheddu <ubaldo@eja.it>

package main

import (
	chatSys "github.com/eja/chat/internal/sys"
	chatWeb "github.com/eja/chat/internal/web"

	"github.com/eja/tibula/log"
	tibulaSys "github.com/eja/tibula/sys"
	tibulaWeb "github.com/eja/tibula/web"
)

func main() {
	if err := chatSys.Configure(); err != nil {
		log.Fatal(err)
	}

	if tibulaSys.Commands.DbSetup {
		if err := tibulaSys.Setup(); err != nil {
			log.Fatal(err)
		}
	} else if tibulaSys.Commands.Wizard {
		if err := tibulaSys.WizardSetup(); err != nil {
			log.Fatal(err)
		}
		if err := chatSys.Wizard(); err != nil {
			log.Fatal(err)
		}

	} else if tibulaSys.Commands.Start {
		if chatSys.Options.DbName == "" {
			log.Fatal("Database name/file is mandatory.")
		}
		if err := chatWeb.Router(); err != nil {
			log.Fatal(err)
		}
		if err := tibulaWeb.Start(); err != nil {
			log.Fatal("Cannot start the web service: ", err)
		}
	} else {
		chatSys.Help()
	}
}
