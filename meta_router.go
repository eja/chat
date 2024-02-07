// Copyright (C) 2023-2024 by Ubaldo Porcheddu <ubaldo@eja.it>

package main

import (
	"net/http"
)

func metaRouter(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		hubMode := r.URL.Query().Get("hub.mode")
		verifyToken := r.URL.Query().Get("hub.verify_token")
		if hubMode == "subscribe" && verifyToken == Options.MetaToken {
			w.Write([]byte(r.URL.Query().Get("hub.challenge")))
		} else {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Token verification error"))
		}
	} else {
		//
	}
}
