// Copyright (C) 2023-2024 by Ubaldo Porcheddu <ubaldo@eja.it>

package main

import (
	"github.com/eja/tibula/db"
)

func translate(languageCode string, label string) string {
	value, err := db.Value("SELECT translation FROM aiTranslations WHERE label=? AND language=? LIMIT 1", label, languageCode)
	if err != nil || value == "" {
		return "{" + label + "}"
	} else {
		return value
	}
}

func languageCodeToLocale(language string) string {
	if locale, err := db.Value("SELECT locale FROM aiLanguages WHERE code = ?", language); err != nil {
		return ""
	} else {
		return locale
	}
}

func languageCodeToInternal(language string) string {
	if value, err := db.Value("SELECT internal FROM aiLanguages WHERE code = ?", language); err != nil {
		return ""
	} else {
		return value
	}
}
