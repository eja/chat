// Copyright (C) 2023-2024 by Ubaldo Porcheddu <ubaldo@eja.it>

package main

import (
	"github.com/eja/chat/internal/core"
	"github.com/eja/chat/internal/router"
	"github.com/eja/tibula/sys"
	"github.com/eja/tibula/web"
	"log"
)

func main() {
	log.SetFlags(0)
	if err := core.Configure(); err != nil {
		log.Fatal(err)
	}

	if sys.Commands.DbSetup {
		log.Println(sys.Options)
		if err := sys.Setup(); err != nil {
			log.Fatal(err)
		}
	} else if sys.Commands.Wizard {
		if err := sys.WizardSetup(); err != nil {
			log.Fatal(err)
		}
		if err := core.Wizard(); err != nil {
			log.Fatal(err)
		}

	} else if sys.Commands.Start {
		if core.Options.DbName == "" {
			log.Fatal("Database name/file is mandatory.")
		}
		if err := router.Router(); err != nil {
			log.Fatal(err)
		}
		if err := web.Start(); err != nil {
			log.Fatal("Cannot start the web service: ", err)
		}
	} else {
		core.Help()
	}
}
