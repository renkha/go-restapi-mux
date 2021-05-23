package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	db "github.com/renkha/go-restapi-gin/config/database"
	"github.com/renkha/go-restapi-gin/middleware"
	"github.com/renkha/go-restapi-gin/src/handler"
	log "github.com/sirupsen/logrus"
)

func init() {
	file, err := os.OpenFile("todo.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println(err)
	}

	log.SetOutput(file)
	log.SetFormatter(&log.TextFormatter{})
	log.SetReportCaller(true)
}

func main() {
	// database
	db := db.DbInit()
	defer db.Close()
	// db.DropTableIfExists(&model.TodoItem{}, &model.User{})
	// db.AutoMigrate(&model.TodoItem{}, &model.User{})

	router := mux.NewRouter()

	// middleware
	router.Use(middleware.Auth(db))

	// try endpoint
	router.HandleFunc("/try", try).Methods("GET")

	// user endpoint
	router.HandleFunc("/api/v1/register", handler.CreateUserHandler(db)).Methods("POST")
	router.HandleFunc("/api/v1/user", handler.GetListUserHandler(db)).Methods("GET")

	// image handler
	router.HandleFunc("/image/{imageName}", handler.ShowImageHandler).Methods("GET")

	// todo endpoint
	router.HandleFunc("/api/v1/todo", handler.CreateToDoHandler(db)).Methods("POST")
	router.HandleFunc("/api/v1/todo", handler.GetListToDoHandler(db)).Methods("GET")
	router.HandleFunc("/api/v1/todo/{id}", handler.GetToDoByIDHandler(db)).Methods("GET")
	router.HandleFunc("/api/v1/todo/{id}", handler.UpdateToDoHandler(db)).Methods("PUT")
	router.HandleFunc("/api/v1/todo/{id}", handler.DeleteToDoHandler(db)).Methods("DELETE")

	http.ListenAndServe(":8080", router)
}

func try(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	io.WriteString(w, `{"alive":true}`)
}
