package model

import "github.com/jinzhu/gorm"

type TodoItem struct {
	gorm.Model
	Description string
	IsCompleted bool
	ImageURL    string
}
