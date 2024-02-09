// Copyright (C) 2023-2024 by Ubaldo Porcheddu <ubaldo@eja.it>

package main

import (
	"fmt"
	"github.com/eja/tibula/db"
)

func dbNumber(value interface{}) int64 {
	return db.Number(value)
}

func dbUserGet(id string) (db.TypeRow, error) {
	return db.Row("SELECT * FROM aiUsers WHERE id = ? AND expiration > CURRENT_TIMESTAMP LIMIT 1", id)
}

func dbUserUpdate(id string, field string, value string) (err error) {
	query := fmt.Sprintf("UPDATE aiUsers SET %s = ? WHERE id = ? LIMIT 1", field)
	_, err = db.Run(query, value, id)
	return
}
