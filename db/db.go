package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"fmt"
)

type Event struct {
	ID int `json:"id" gorm:"id"`
	Name string `json:"name" gorm:"name"`
	Data string `json:"data" gorm:"data"`
	CreatedAt int64 `json:"created_at" gorm:"created_at"`
	UpdatedAt int64 `json:"updated_at" gorm:"update_at"`
}

type User struct {
	//ID int `json:"id" gorm:"id"`
	Email string `json:"email" gorm:"email"`
	Password string `json:"password" gorm:"password"`
}

func DatabaseConnect() *gorm.DB {
	database, err :=
		gorm.Open("postgres", "host=localhost user=root dbname=test sslmode=disable password=123")
	if err != nil {
		fmt.Printf("Panic! At the disco.")
	}
	//defer database.Close()
	return database
}
