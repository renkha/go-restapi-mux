package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/renkha/go-restapi-gin/src/model"
	log "github.com/sirupsen/logrus"
)

func CreateToDoHandler(db *gorm.DB) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		description := r.FormValue("description")

		fileUploaded, header, err := r.FormFile("image")
		if err != nil {
			log.Warn("Failed get image with error" + err.Error())
			w.Header().Set("content-type", "applocation/json")
			w.WriteHeader(http.StatusBadRequest)
			io.WriteString(w, `{"Error":"Failed to get image"}`)
			return
		}
		defer fileUploaded.Close()

		dir, err := os.Getwd()
		if err != nil {
			log.Warn("Failed get image working directory" + err.Error())
			w.Header().Set("content-type", "applocation/json")
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, `{"Error":"Failed to get wworking directory"}`)
			return
		}

		t := time.Now()

		fileName := fmt.Sprintf("%s%s%s", "image", t.Format("20060102150405"), filepath.Ext(header.Filename))
		fileLocation := filepath.Join(dir, "images", fileName)
		targetFile, err := os.OpenFile(fileLocation, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			log.Warn("Failed get file with error" + err.Error())
			w.Header().Set("content-type", "applocation/json")
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, `{"Error":"Failed to open file"}`)
			return
		}
		defer targetFile.Close()

		_, err = io.Copy(targetFile, fileUploaded)
		if err != nil {
			log.Warn("Failed copy file with error" + err.Error())
			w.Header().Set("content-type", "applocation/json")
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, `{"Error":"Failed to copy file"}`)
			return
		}

		fullPathName := fmt.Sprintf("%s/image/%s", r.Host, fileName)

		newToDo := &model.TodoItem{
			Description: description,
			IsCompleted: false,
			ImageURL:    fileName,
		}

		db.Create(&newToDo)
		result := db.Last(&newToDo)
		log.WithFields(log.Fields{"Description": description, "ImageURL": fullPathName}).Info("Succes to add image")

		w.Header().Set("content-type", "application/json")
		json.NewEncoder(w).Encode(result.Value)
	}

	return http.HandlerFunc(fn)
}

func GetListToDoHandler(db *gorm.DB) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		todoList := []model.TodoItem{}

		db.Find(&todoList)

		w.Header().Set("content-type", "application/json")
		json.NewEncoder(w).Encode(todoList)
	}

	return http.HandlerFunc(fn)
}

func GetToDoByIDHandler(db *gorm.DB) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var todo model.TodoItem
		vars := mux.Vars(r)
		id, _ := strconv.Atoi(vars["id"])

		db.Where("id = ? ", id).First(&todo)

		w.Header().Set("content-type", "application/json")
		json.NewEncoder(w).Encode(todo)
	}

	return http.HandlerFunc(fn)
}

func UpdateToDoHandler(db *gorm.DB) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var todo = model.TodoItem{}
		vars := mux.Vars(r)
		id, _ := strconv.Atoi(vars["id"])

		description := r.FormValue("description")
		isCompleted, _ := strconv.ParseBool(r.FormValue("is_completed"))

		db.Where("id = ? ", id).First(&todo)
		todo.Description = description
		todo.IsCompleted = isCompleted

		db.Save(&todo)
		log.WithFields(log.Fields{
			"ID":          id,
			"Description": description,
			"IsCompleted": isCompleted,
		}).Info("Succes updating todo item")

		w.Header().Set("content-type", "application/json")
		json.NewEncoder(w).Encode(todo)
	}

	return http.HandlerFunc(fn)
}

func DeleteToDoHandler(db *gorm.DB) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var todo = model.TodoItem{}
		vars := mux.Vars(r)
		id, _ := strconv.Atoi(vars["id"])

		db.Where("id = ? ", id).Delete(&todo)

		log.WithFields(log.Fields{
			"ID": id,
		}).Info("Success delete todo item")

		w.Header().Set("content-type", "application/json")
		io.WriteString(w, `{"success":true}`)
	}

	return http.HandlerFunc(fn)
}
