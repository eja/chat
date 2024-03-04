// Copyright (C) 2023-2024 by Ubaldo Porcheddu <ubaldo@eja.it>

package web

import (
	"fmt"

	"github.com/eja/chat/internal/process"
	"github.com/eja/tibula/api"
	"github.com/eja/tibula/web"
)

func Router() error {
	web.Router.HandleFunc("/meta", metaRouter)
	web.Router.HandleFunc("/tg", telegramRouter)

	api.Plugins["aiChat"] = func(eja api.TypeApi) api.TypeApi {
		if eja.Action == "run" && eja.Values["chat"] != "" {
			user := fmt.Sprintf("T.%d", eja.Owner)
			language := eja.Language
			if answer, err := process.Text(user, language, eja.Values["chat"]); err != nil {
				eja.Alert = append(eja.Alert, err.Error())
			} else {
				eja.Values["chat"] = answer
			}
		}
		return eja
	}
	return nil
}
