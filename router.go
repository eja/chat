// Copyright (C) 2023-2024 by Ubaldo Porcheddu <ubaldo@eja.it>

package main

import (
	"github.com/eja/tibula/web"
)

func Router() error {
	web.Router.HandleFunc("/meta", metaRouter)
	web.Router.HandleFunc("/tg", telegramRouter)

	return nil
}
