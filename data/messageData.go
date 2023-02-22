package data

import (
	"gorm.io/gorm"
	"time"
)

type DouyinMessageActionRequest struct {
	Content  string // 消息内容
	UserId   int64  // 用户Id，token鉴权得到
	ToUserId int64
}
type DouyinMessageHistoryRequest struct {
	UserId   int64 // 用户Id，token鉴权得到
	ToUserId int64
	LastTime time.Time
}
type DouyinMessageHistoryResponse struct {
	Response
	Messages []*Message
}
type Message struct {
	Id         int64     `gorm:"primarykey" json:"id,omitempty"`
	UserID     int64     `gorm:"column:user_id" json:"user_id,omitempty"`
	ToUserID   int64     `gorm:"column:to_user_id" json:"to_user_id,omitempty"`
	Content    string    `gorm:"column:content" json:"content,omitempty"`
	CreateTime time.Time `gorm:"column:created_at" json:"create_time"`
}
type MessageRaw struct {
	gorm.Model
	UserID   int64  `gorm:"column:user_id"`
	ToUserID int64  `gorm:"column:to_user_id"`
	Content  string `gorm:"column:content"`
}

func (MessageRaw) TableName() string {
	return "messages"
}
func (Message) TableName() string {
	return "messages"
}
