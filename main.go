// Copyright (C) 2023-2024 by Ubaldo Porcheddu <ubaldo@eja.it>

package main

import (
	"github.com/eja/tibula/sys"
	"github.com/eja/tibula/web"
	"log"
)

func main() {
	log.SetFlags(0)
	if err := Configure(); err != nil {
		log.Fatal(err)
	}
	if sys.Commands.Start {
		if Options.DbName == "" {
			log.Fatal("Database name/file is mandatory.")
		}
		if err := Router(); err != nil {
			log.Fatal(err)
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
