package service

import (
	"dousheng/data"
)

type CommentService interface {
	CommentAction(user_id int64, vidio_id int64, action_type int32, comment_text string, comment_id int64) (data.Comment, error)
	CommentList(video_id int64) ([]data.Comment, error)
}

