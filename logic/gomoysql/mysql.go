package gomoysql

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

type User struct {
	Username string
	Password string
	ID       uint `gorm:"primarykey"`
}

var Db *gorm.DB

func InitMysql() {
	dsn := "root:zhs123..@tcp(192.168.10.128:3306)/gochat?charset=utf8mb4"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("mysql connect err", err)
	}
	Db = db
	err = Db.AutoMigrate(&User{})
	if err != nil {
		log.Fatal("create table err", err)
	}
}

func AddUser(user *User) {
	Db.Create(user)
}

func QueryUsername(username string) User {
	var user User
	Db.Where("username = ?", username).First(&user)
	return user
}
