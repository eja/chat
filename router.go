// Copyright (C) 2023-2024 by Ubaldo Porcheddu <ubaldo@eja.it>

package main

import (
	"github.com/eja/tibula/db"
	"github.com/eja/tibula/sys"
	"github.com/eja/tibula/web"
)

func Router() error {
	db.LogLevel = sys.Options.LogLevel
	if err := db.Open(sys.Options.DbType, sys.Options.DbName, sys.Options.DbUser, sys.Options.DbPass, sys.Options.DbHost, sys.Options.DbPort); err != nil {
		return err
	}

	web.Router.HandleFunc("/meta", metaRouter)
	web.Router.HandleFunc("/tg", telegramRouter)

	return nil
}
