package handler

import (
	"encoding/json"
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/renkha/go-restapi-mux/src/model"
)

func CreateUserHandler(db *gorm.DB) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		name := r.FormValue("name")
		pass := r.FormValue("password")

		newUser := &model.User{
			Name:     name,
			Password: pass,
		}

		db.Create(&newUser)
		result := db.Last(&newUser)

		w.Header().Set("content-type", "application/json")
		json.NewEncoder(w).Encode(result.Value)
	}

	return http.HandlerFunc(fn)
}

func GetListUserHandler(db *gorm.DB) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		users := []model.User{}

		db.Find(&users)

		w.Header().Set("content-type", "application/json")
		json.NewEncoder(w).Encode(users)
	}

	return http.HandlerFunc(fn)
}
