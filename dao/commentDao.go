package dao

import "dousheng/data"

// import "dousheng/data"

func CommentCount() (int64, error) {
	db := GetDB()
	var count int64
	res := db.Table("comments").Select("MAX(id)").Find(&count)
	return count, res.Error
}

func InsertComment(user_id int64, video_id int64, commentText string) error {
	db := GetDB()
	comment := comments{
		UserId:      user_id,
		VideoId:     video_id,
		CommentText: commentText,
	}
	res := db.Table("comments").Select("user_id", "video_id", "comment_text").Create(&comment)
	
	return res.Error
}

func DeleteteComment(comment_id int64) error {
	result := DB.Model(comments{}).Where(map[string]interface{}{"id": comment_id}).
		Update("is_comment", 0)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func GetCommentList(video_id int64) ([]data.DBComments, error) {
	db := GetDB()
	dbcomments := make([]data.DBComments, 0)
	res := db.Table("comments").Where("video_id = ? AND is_comment = 1", video_id).Find(&dbcomments)
	return dbcomments, res.Error
}
