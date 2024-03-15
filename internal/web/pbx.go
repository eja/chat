// Copyright (C) 2023-2024 by Ubaldo Porcheddu <ubaldo@eja.it>

package web

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/eja/chat/internal/db"
	"github.com/eja/chat/internal/process"
	"github.com/eja/chat/internal/sys"
	"github.com/eja/tibula/log"
)

const pbxMaxMemory = 32 * 1024 * 1024 //32MB

func pbxRouter(w http.ResponseWriter, r *http.Request) {
	if sys.Options.PbxToken == "" {
		w.WriteHeader(http.StatusNotImplemented)
		return
	}

	if r.Method == "POST" {
		if err := db.Open(); err != nil {
			return
		}

		now := strconv.FormatInt(time.Now().Unix(), 10)

		err := r.ParseMultipartForm(pbxMaxMemory)
		if err != nil {
			errMessage := "Error parsing form-data"
			http.Error(w, errMessage, http.StatusBadRequest)
			log.Warn("[PBX]", errMessage)
			return
		}

		callerId := r.FormValue("phone")
		authToken := r.FormValue("token")
		fileInput, _, err := r.FormFile("file")
		if err != nil {
			errMessage := "Error retrieving file from form-data"
			http.Error(w, errMessage, http.StatusBadRequest)
			log.Warn("[PBX]", errMessage)
			return
		}
		defer fileInput.Close()

		fileInputName := fmt.Sprintf("%s.%s.pbx.in", callerId, now)

		if authToken != sys.Options.PbxToken {
			errMessage := "Invalid authorization token"
			http.Error(w, errMessage, http.StatusUnauthorized)
			log.Debug("[PBX]", errMessage)
			return
		}

		user, err := db.UserGet(callerId)
		if err != nil || user == nil {
			errMessage := "User not found"
			http.Error(w, errMessage, http.StatusUnauthorized)
			log.Debug("[PBX]", errMessage)
			return
		}

		out, err := os.Create(fmt.Sprintf("%s/%s", sys.Options.MediaPath, fileInputName))
		if err != nil {
			errMessage := "Error creating file on server"
			http.Error(w, errMessage, http.StatusInternalServerError)
			log.Warn("[PBX]", errMessage)
			return
		}
		defer out.Close()

		_, err = io.Copy(out, fileInput)
		if err != nil {
			errMessage := "Error copying file on server"
			http.Error(w, errMessage, http.StatusInternalServerError)
			log.Warn("[PBX]", errMessage)
			return
		}

		fileOutputName, err := process.Audio("pbx", user["id"], user["language"], "", fileInputName, true)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Warn("[PBX]", err.Error())
			return
		}
		if fileOutputName == "" {
			errMessage := "Empty response"
			http.Error(w, errMessage, http.StatusInternalServerError)
			log.Warn("[PBX]", errMessage)
			return
		}

		fileOutput, err := os.Open(fileOutputName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Warn("[PBX]", err.Error())
			return
		}
		defer fileOutput.Close()

		w.Header().Set("Content-Disposition", "attachment; filename="+fileOutputName)
		if _, err := io.Copy(w, fileOutput); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Warn("[PBX]", err.Error())
			return
		}

	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
