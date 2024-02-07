// Copyright (C) 2023-2024 by Ubaldo Porcheddu <ubaldo@eja.it>

package main

import (
	"github.com/eja/tibula/db"
)

func Translate(languageCode string, label string) string {
	return "{" + label + "}"
}

func LanguageCodeToLocale(language string) string {
	if locale, err := db.Value("SELECT locale FROM aiLanguages WHERE code = ?", language); err != nil {
		return ""
	} else {
		return locale
	}
}

func LanguageCodeToInternal(language string) string {
	if value, err := db.Value("SELECT internal FROM aiLanguages WHERE code = ?", language); err != nil {
		return ""
	} else {
		return value
	}
}
