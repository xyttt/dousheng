package dao

import (
	"dousheng/config"
	"dousheng/data"
	"time"
)

func SelectFeedByTimeId(latestTime time.Time, userID int64) ([]data.DBVideo, error) {
	db := GetDB()
	dbvideos := make([]data.DBVideo, config.VideoNum)
	res := db.Table("videos").Where("publish_time < ? AND author_id != ?", latestTime, userID).Order("publish_time desc").Limit(config.VideoNum).Find(&dbvideos)
	return dbvideos, res.Error
}

func CountFavorite(videoId int64) (int64, error) {
	db := GetDB()
	var count int64
	res := db.Table("favorites").Where("video_id = ? AND action_type = ?", videoId, 1).Count(&count)
	return count, res.Error
}

func CountComment(videoId int64) (int64, error) {
	db := GetDB()
	var count int64
	res := db.Table("comments").Where("video_id = ? AND is_comment = 1", videoId).Count(&count)
	return count, res.Error
}

func JudgeIsFavorite(videoId int64, userID int64) (bool, error) {
	db := GetDB()
	var count int64
	//TODO：可改进
	res := db.Table("favorites").Where("video_id = ? AND user_id = ? AND action_type = 1", videoId, userID).Count(&count)
	if count > 0 {
		return true, res.Error
	} else {
		return false, res.Error
	}
}

func SelectNameByUserId(authorId int64) (string, error) {
	db := GetDB()
	var s string
	res := db.Table("users").Select("name").Where("id = ?", authorId).Find(&s)
	return s, res.Error
}

// 关注
func CountFollow(authorId int64) (int64, error) {
	db := GetDB()
	var count int64
	res := db.Table("follows").Where("user_id = ? AND is_follow = 1", authorId).Count(&count)
	return count, res.Error
}

// 粉丝
func CountFollower(authorId int64) (int64, error) {
	db := GetDB()
	var count int64
	res := db.Table("follows").Where("followed_id = ? AND is_follow = 1", authorId).Count(&count)
	return count, res.Error
}

func JudgeIsFollow(authorId int64, userId int64) (bool, error) {
	db := GetDB()
	var count int64
	//TODO：可改进
	res := db.Table("follows").Where("user_id = ? AND followed_id = ? AND is_follow = 1", userId, authorId).Count(&count)
	if count > 0 {
		return true, res.Error
	} else {
		return false, res.Error
	}
}

func SelectVideoByUserId(userID int64) ([]data.DBVideo, error) {
	db := GetDB()
	dbvideos := make([]data.DBVideo, 0)
	res := db.Table("videos").Where("author_id = ?", userID).Order("publish_time desc").Find(&dbvideos)
	return dbvideos, res.Error
}

func InsertVideo(authorId int64, playUrl string, coverUrl string, title string) error {
	db := GetDB()
	video := data.DBVideo{
		AuthorId: authorId,
		PlayUrl:  playUrl,
		CoverUrl: coverUrl,
		Title:    title,
	}
	res := db.Table("videos").Select("author_id", "play_url", "cover_url", "title").Create(&video)
	return res.Error
}
