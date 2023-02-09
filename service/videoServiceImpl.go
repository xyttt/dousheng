package service

import (
	"dousheng/config"
	"dousheng/dao"
	"dousheng/data"
	"log"
	"time"
)

// 组装user
func PkgUserByAuthorid(authorId int64, userId int64, user *data.User) error {
	//var user data.User
	user.Id = authorId
	var err error
	user.Name, err = dao.SelectNameByUserId(authorId)
	user.FollowCount, err = dao.CountFollow(authorId)
	user.FollowerCount, err = dao.CountFollower(authorId)
	user.IsFollow, err = dao.JudgeIsFollow(authorId, userId)
	return err
}

func PkgVideosByVideoid(userId int64, dbvideos []data.DBVideo, videos *[]data.Video) error {
	var err error
	for _, dbvideo := range dbvideos {
		var video data.Video
		video.DBVideo = dbvideo

		err = PkgUserByAuthorid(dbvideo.AuthorId, userId, &video.Author)
		if err != nil {
			log.Printf("failed with GetUserByAuthorid(video.AuthorId): %v", err)
		}
		video.FavoriteCount, err = dao.CountFavorite(dbvideo.Id)
		if err != nil {
			log.Printf("failed with dao.CountFavorite(video.Id): %v", err)
		}
		video.CommentCount, err = dao.CountComment(dbvideo.Id)
		if err != nil {
			log.Printf("failed with dao.CountComment(video.Id): %v", err)
		}
		video.IsFavorite, err = dao.JudgeIsFavorite(dbvideo.Id, userId)
		if err != nil {
			log.Printf("failed with dao.CountIsFavorite(video.Id, userID): %v", err)
		}
		*videos = append(*videos, video)

	}
	return err
}

func (v VideoServiceImpl) Feed(latestTime time.Time, userId int64) ([]data.Video, int64, error) {
	videos := make([]data.Video, 0, config.VideoNum)
	//dbvideos := make([]data.DBVideo, config.VideoNum)

	dbvideos, err := dao.SelectFeedByTimeId(latestTime, userId)
	if err != nil {
		log.Printf("failed with dao.SelectFeedByTimeId(latestTime, userID): %v", err)
		return nil, time.Now().Unix(), err
	}

	err = PkgVideosByVideoid(userId, dbvideos, &videos)
	if err != nil {
		log.Printf("failed with PkgVideosByVideoid: %v", err)
		return nil, time.Now().Unix(), err
	}

	nestTime := dbvideos[len(dbvideos)-1].PublishTime
	return videos, nestTime.Unix(), err
}

func (v VideoServiceImpl) PubList(userId int64) ([]data.Video, error) {
	videos := make([]data.Video, 0)

	dbvideos, err := dao.SelectVideoByUserId(userId)
	if err != nil {
		log.Printf("failed with dao.SelectVideoByUserId(userId): %v", err)
		return nil, err
	}

	err = PkgVideosByVideoid(userId, dbvideos, &videos)
	if err != nil {
		log.Printf("failed with PkgVideosByVideoid: %v", err)
		return nil, err
	}

	return videos, err

}
