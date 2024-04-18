package main

import (
	"dy/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(mysql.Open("root:root@tcp(111.230.33.57:3306)/go_chat?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	//迁移schema
	db.AutoMigrate(&models.Contact{})

	//msg := &models.GroupBasic{}
	//msg.FromId = 1
	//msg.TargetId = 2
	//db.Create(msg)
	//
	//fmt.Println(db.First(msg, 1))
	//
	//db.Model(msg).Update("target_id", "4")
}
