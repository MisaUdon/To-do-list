package db

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Todo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

// var Todos []Todo

func Init() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&Todo{})
}

func CreateTodo(todo Todo) error {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
		return err
	}
	res := db.Create(&todo)
	return res.Error
}

func GetAllTodos() ([]Todo, error) {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	var todos []Todo
	res := db.Find(&todos)
	return todos, res.Error
}

func UpdateTodo(todo Todo) error {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
		return err
	}
	res := db.Save(&todo)
	return res.Error
}

func DeleteTodo(id string) error {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
		return err
	}
	res := db.Delete(&Todo{}, "id = ?", id)
	return res.Error
}
