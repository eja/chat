// Copyright (C) 2023-2024 by Ubaldo Porcheddu <ubaldo@eja.it>

package main

import (
	"github.com/eja/tibula/db"
	"github.com/eja/tibula/sys"
)

func defaultLanguage() string {
	language, err := db.Value("SELECT code FROM aiLanguages WHERE default_language > 0")
	if err != nil || language == "" {
		language = sys.Options.Language
	}
	return language
}

func translate(language string, label string) string {
	if language == "" {
		language = defaultLanguage()
	}
	value, err := db.Value("SELECT translation FROM aiTranslations WHERE label=? AND language=?", label, language)
	if err != nil || value == "" {
		return "{" + label + "}"
	} else {
		return value
	}
}

func languageCodeToLocale(language string) string {
	if language == "" {
		language = defaultLanguage()
	}

	if locale, err := db.Value("SELECT locale FROM aiLanguages WHERE code = ?", language); err != nil {
		return ""
	} else {
		return locale
	}
}

func languageCodeToInternal(language string) string {
	if language == "" {
		language = defaultLanguage()
	}

	if value, err := db.Value("SELECT internal FROM aiLanguages WHERE code = ?", language); err != nil {
		return ""
	} else {
		return value
	}
}
