package db

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID       uint   `gorm:"primary_key"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

var DBASE *gorm.DB

func Connect() *gorm.DB {
	dbURL := "host=localhost user=postgres password=vamsi123 dbname=soulpage port=5432 sslmode=disable"
	Database, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true})
	if err != nil {
		fmt.Println(err.Error())
		panic("Cannot connect to DB")
	} else {
		fmt.Println("connected to the database successfully")
	}

	err = Database.AutoMigrate(&User{})
	if err != nil {
		fmt.Println(err.Error())
	}

	DBASE = Database
	return DBASE
}
