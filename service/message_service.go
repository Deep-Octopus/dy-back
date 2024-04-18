package service

import "dy/models"

func GetMessagesByContact(contact models.Contact) []models.Message {
	msgs := make([]models.Message, 0)
	//发送的消息
	msgs = append(msgs, models.GetMessagesByFromIdAndTargetIdAndType(contact.OwnerId, contact.TargetId, uint(contact.Type))...)
	//收到的消息
	msgs = append(msgs, models.GetMessagesByFromIdAndTargetIdAndType(contact.TargetId, contact.OwnerId, uint(contact.Type))...)
	return msgs
}

func GetGroupMessagesByContact(contact models.Contact) []models.Message {
	return models.GetGroupMessagesByTargetIdAndType(contact.TargetId, uint(contact.Type))
}
