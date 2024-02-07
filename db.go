// Copyright (C) 2023-2024 by Ubaldo Porcheddu <ubaldo@eja.it>

package main

import (
	"fmt"
	"github.com/eja/tibula/db"
)

func DbUserGet(id string) (db.TypeRow, error) {
	return db.Row("SELECT * FROM aiUsers WHERE id = ? AND expiration > CURRENT_TIMESTAMP LIMIT 1", id)
}

func DbUserUpdate(id string, field string, value string) (err error) {
	_, err = db.Run(fmt.Sprintf("UPDATE aiUsers SET %s = ? WHERE id = ? LIMIT 1", field), id, value)
	return
}
