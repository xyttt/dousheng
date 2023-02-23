package service

import (
	"dousheng/dao"
	"dousheng/data"
	"time"
)

type CommentServiceImpl struct {
}

func (coni *CommentServiceImpl) CommentAction(user_id int64, vidio_id int64, action_type int32, comment_text string, comment_id int64) (data.Comment, error) {
	if action_type == 1 {
		dao.InsertComment(user_id, vidio_id, comment_text)
	}
	id, _ := dao.CommentCount()
	if action_type == 2 {
		dao.DeleteteComment(comment_id)
		id = comment_id
	}

	now := time.Now()
	strTime := now.Format("2006-01-02 15:04:05")
	u, _ := dao.GetTableUserById(user_id)
	user := data.User{
		Id:   u.Id,
		Name: u.Name,
	}

	result := data.Comment{
		Id:          id,
		User:        user,
		Content:     comment_text,
		Create_date: strTime,
	}
	return result, nil
}

func (coni *CommentServiceImpl) CommentList(video_id int64) ([]data.Comment, error) {
	dbcomments,_ := dao.GetCommentList(video_id)
	coms := make([]data.Comment, 0)
	CommentList_(dbcomments, &coms)

	return coms,nil
}

func CommentList_(dbcomments []data.DBComments, comments *[]data.Comment) error {
	var err error
	for _, dbcomment := range dbcomments {
		now := time.Now()
		strTime := now.Format("2006-01-02 15:04:05")
		u, _ := dao.GetTableUserById(dbcomment.UserId)
		user := data.User{
			Id:   u.Id,
			Name: u.Name,
		}

		comment := data.Comment{
			Id:          dbcomment.Id,
			User:        user,
			Content:     dbcomment.CommentText,
			Create_date: strTime,
		}

		*comments = append(*comments, comment)
	}
	return err
}