package dao

import (
	"dousheng/data"
)

func CreateMes(UserID int64, To_UserId int64, Content string) error {
	err := DB.Create(&data.Message{UserID: UserID, ToUserID: To_UserId, Content: Content}).Error
	return err
}

func MessageHistory(UserID int64, To_UserId int64, LastTime int64) ([]*data.Message, error) {
	var Message []*data.Message
	//TODO：加入时间条件
	err := DB.Debug().Where("user_id = ? and to_user_id = ?  or to_user_id = ? and user_id = ?",
		UserID, To_UserId, UserID, To_UserId).Order("created_at").Find(Message).Error
	if err != nil {
		return nil, err
	}
	return Message, nil
}
