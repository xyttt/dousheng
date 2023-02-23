package data

import (
	"time"
)

type Comment struct {
	Id             int64 		`json:"id"`
	User       	   User			`json:"user"`
	Content		   string 		`json:"content"`
	Create_date    string  		`json:"create_date"`
}

type DBComments struct {
	Id          int64     `json:"id"`
	UserId      int64     `json:"Userid"`
	VideoId     int64     `json:"video_id"`
	CommentText string    `json:"comment_text"`
	PublishTime time.Time `json:"publish_time"`
	IsComment   int8      `json:"id_comment"`
}

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

type Empty struct {
}
