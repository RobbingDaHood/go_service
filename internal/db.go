package internal

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"todo-service/internal/model"
)

var DB *gorm.DB

func Init() {
	var err error
	DB, err = gorm.Open(sqlite.Open("todo.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Auto-migrate the Todo model
	DB.AutoMigrate(&model.Todo{})
}
