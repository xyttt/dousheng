package data

import "time"

type DouyinMessageActionRequest struct {
	Content  string // 消息内容
	UserId   int64  // 用户Id，token鉴权得到
	ToUserId int64
}
type DouyinMessageHistoryRequest struct {
	UserId   int64 // 用户Id，token鉴权得到
	ToUserId int64
	LastTime int64
}
type DouyinMessageHistoryResponse struct {
	Response
	Messages []*Message
}
type Message struct {
	Id         int64
	UserID     int64     `gorm:"user_id"`
	ToUserID   int64     `gorm:"to_user_id"`
	Content    string    `gorm:"content"`
	CreateTime time.Time `gorm:"created_at"`
}

func (Message) TableName() string {
	return "messages"
}
