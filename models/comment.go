package models

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	UserID  uint   `json:"userID"`
	VideoID uint   `json:"videoID"`
	Text    string `json:"text"`
	// 其他评论信息字段
}
