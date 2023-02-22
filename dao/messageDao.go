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
	//TODO：加入时间条件
	err := DB.Debug().Where("created_at < ?", LastTime).Where("user_id = ? and to_user_id = ?  or to_user_id = ? and user_id = ?",
		UserID, To_UserId, UserID, To_UserId).Order("created_at").Find(&Message).Error
	//这里是否考虑联合索引？
	if err != nil {
		return nil, err
	}
	return Message, nil
}
