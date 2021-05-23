package handler

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func ShowImageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	imgName := vars["imageName"]

	fileDir := fmt.Sprintf("./images/%s", imgName)
	file, err := os.Open(fileDir)
	if err != nil {
		log.Warn("Failed get file with error" + err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer file.Close()

	io.Copy(w, file)
}
