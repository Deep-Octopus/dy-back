package sql

import (
	"dy/models"
	"dy/utils"
	"fmt"
)

func InitMysqlSchema() {
	//迁移schema
	err := utils.DB.AutoMigrate(&models.Contact{},
		&models.Message{},
		&models.UserBasic{},
		&models.GroupBasic{},
		&models.Comment{},
		&models.Video{},
		&models.UserAndVideo{})
	if err != nil {
		return
	}
	fmt.Println("迁移数据库Schema成功")
}
