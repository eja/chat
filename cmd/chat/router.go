// Copyright (C) 2023-2024 by Ubaldo Porcheddu <ubaldo@eja.it>

package main

import (
	"github.com/eja/tibula/web"
)

func Router() error {
	web.Router.HandleFunc("/meta", MetaRouter)
	web.Router.HandleFunc("/tg", TelegramRouter)
	return nil
}
