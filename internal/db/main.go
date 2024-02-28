// Copyright (C) 2023-2024 by Ubaldo Porcheddu <ubaldo@eja.it>

package db

import (
	"fmt"

	"github.com/eja/chat/internal/sys"
	"github.com/eja/tibula/db"
)

func Number(value interface{}) int64 {
	return db.Number(value)
}

func UserGet(id string) (db.TypeRow, error) {
	return db.Row("SELECT * FROM aiUsers WHERE id = ? AND expiration > CURRENT_TIMESTAMP LIMIT 1", id)
}

func UserUpdate(id string, field string, value string) (err error) {
	query := fmt.Sprintf("UPDATE aiUsers SET %s = ? WHERE id = ?", field)
	_, err = db.Run(query, value, id)
	return
}

func SystemPrompt() (db.TypeRows, error) {
	return db.Rows("SELECT prompt FROM aiPrompts WHERE active > 0 ORDER BY power ASC")
}

func Open() error {
	if err := db.Open(sys.Options.DbType, sys.Options.DbName, sys.Options.DbUser, sys.Options.DbPass, sys.Options.DbHost, sys.Options.DbPort); err != nil {
		return err
	}
	return nil
}
