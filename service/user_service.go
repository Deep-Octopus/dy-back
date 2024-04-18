package service

import (
	"dy/models"
	"dy/utils"
	"errors"
	"github.com/asaskevich/govalidator"
	"time"
)

func GetUserList() []*models.UserBasic {
	return models.GetUserList()
}

func GetUserByUsername(username string) *models.UserBasic {
	return models.TakeUserByUsername(username)
}
func GetUserByEmail(email string) *models.UserBasic {
	return models.TakeUserByEmail(email)
}
func GetFriendByUsername(userId uint) []models.UserBasic {
	return models.SearchFriend(userId)
}
func GetGroupByUsername(userId uint) []models.GroupBasic {
	return models.SearchGroup(userId)
}
func GetListMessageByUsername(userId uint) []models.ListMessageDto {
	lms := make([]models.ListMessageDto, 0)
	userLms := models.SearchFriendAndLastMessage(userId)
	groupLms := models.SearchGroupAndLastMessage(userId)
	lms = append(lms, groupLms...)
	lms = append(lms, userLms...)
	return lms
}

func CreateUser(user *models.UserBasic) error {
	if models.TakeUserByEmail(user.Email).ID != 0 {
		return errors.New("用户已经存在")
	}
	user.LogoutTime = time.Now()
	user.LoginTime = time.Now()
	user.HeartbeatTime = time.Now()
	user.Username = utils.GenerateAccountNumber(10)
	user.Password, _ = utils.HashPassword(user.Password)
	if len(user.Name) == 0 {
		user.Name = "抖音用户@" + user.Username
	}
	if err := models.CreateUser(user); err != nil {
		return errors.New("新增用户失败")
	}
	return nil
}

func DeleteUser(user *models.UserBasic) error {
	if err := models.DeleteUser(user); err != nil {
		return errors.New("删除用户失败")
	}
	return nil
}

func UpdateUser(user *models.UserBasic) error {
	if _, err := govalidator.ValidateStruct(user); err != nil {
		return errors.New("数据格式错误")
	}
	if user.Password != "" {
		user.Password, _ = utils.HashPassword(user.Password)
	}
	if err := models.UpdateUser(user); err != nil {
		return errors.New("更新用户失败")
	}
	return nil
}
