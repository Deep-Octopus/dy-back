package models

import (
	"dy/utils"
	"gorm.io/gorm"
)

// GroupBasic 群聊
type GroupBasic struct {
	gorm.Model
	Name    string //群聊名称
	OwnerId uint   //群聊拥有者
	Icon    string //图标
	Type    int
	Desc    string // 预留
}

func GetGroupMembersByGroupID(groupId uint) []int64 {
	var ownerIds []int64
	utils.DB.Where("target_id = ? and type = 2", groupId).Find(&ownerIds, "owner_id")
	return ownerIds
}
