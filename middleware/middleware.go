package middleware

import (
	"io"
	"net/http"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/renkha/go-restapi-mux/src/model"
	log "github.com/sirupsen/logrus"
)

func Auth(db *gorm.DB) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Info("Access from ", r.URL.Path)

			notAuth := []string{
				"/try",
				"/api/v1/register",
				"/image",
			}
			withAuth := true
			for _, path := range notAuth {
				if path == r.URL.Path || strings.HasPrefix(r.URL.Path, path) {
					withAuth = false
					break
				}
			}
			name, password, ok := r.BasicAuth()

			if withAuth {
				if !ok {
					w.Header().Set("content-type", "application/json")
					w.WriteHeader(http.StatusInternalServerError)
					io.WriteString(w, `{"Error":"Internal Server Error"}`)
					return
				}

				var user model.User
				tx := db.Where("name = ? AND password = ?", name, password).First(&user)
				if tx.Error != nil || tx.RowsAffected == 0 {
					w.Header().Set("content-type", "application/json")
					w.WriteHeader(http.StatusUnauthorized)
					io.WriteString(w, `{"Error":"User not Authorized"}`)
					return
				}
			}

			h.ServeHTTP(w, r)
		})
	}
}
