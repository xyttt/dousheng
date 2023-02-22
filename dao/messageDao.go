package dao

import (
	"dousheng/data"
	"time"
)

func CreateMes(UserID int64, To_UserId int64, Content string) error {
	err := DB.Create(&data.MessageRaw{UserID: UserID, ToUserID: To_UserId, Content: Content}).Error
	return err
}

func MessageHistory(UserID int64, To_UserId int64, LastTime time.Time) ([]*data.Message, error) {
	var Message []*data.Message
	err := DB.Debug().Where("created_at > ?", LastTime).Where("user_id = ? and to_user_id = ?  or user_id = ? and to_user_id = ?",
		UserID, To_UserId, To_UserId, UserID).Order("created_at").Find(&Message).Error
	if err != nil {
		return nil, err
	}
	return Message, nil
}
