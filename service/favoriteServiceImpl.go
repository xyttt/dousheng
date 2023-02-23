package service

import (
	"dousheng/dao"
	"dousheng/data"
	"log"
)

type FavoriteServiceImpl struct {
}

func (favorite *FavoriteServiceImpl) FavoriteAction(user_id int64, vidio_id int64, action_type int32) error {
	// 查询是否有点赞记录
	favoriteData, err := dao.GetFavorite(user_id, vidio_id)
	if err != nil {
		return err
	}
	// 如果没有，就插入点赞数据
	if favoriteData.Id == 0 {
		dao.InsertFavorite(user_id, vidio_id)
		return nil
	}
	// 如果有，则根据action_type执行操作
	dao.UpdateFavorite(user_id, vidio_id, action_type)
	return nil
}

func (favorite *FavoriteServiceImpl) FavoriteList(user_id int64, token string) ([]data.FavoriteVideo, error) {
	favoriteList, err1 := dao.GetFavoriteList(user_id)
	log.Println("favoritelist")
	log.Println(favoriteList)
	if err1 != nil {
		log.Printf("方法:GetFavouriteList query key失败: %v", err1)
		return nil, err1
	}
	favoriteVideoList, err2 := dao.GetVideoListByVideoIds(favoriteList)
	log.Println(favoriteVideoList)
	if err2 != nil {
		log.Printf("方法:GetVideoListByVideoIds query key失败: %v", err2)
		return nil, err2
	}
	videos, err3 := favorite.PkgFavoriteList(favoriteVideoList)
	if err3 != nil {
		log.Printf("方法:GPkgFavoriteList失败: %v", err2)
		return nil, err3
	}
	return videos, nil
}

// 将DBVideo组装FavoriteVideo
func (favorite *FavoriteServiceImpl) PkgFavoriteList(favoriteVideoList []data.Video) ([]data.FavoriteVideo, error) {
	var videos []data.FavoriteVideo
	for _, dbvideo := range favoriteVideoList {
		var video data.FavoriteVideo
		err := PkgVideoByVideoid(dbvideo.AuthorId, dbvideo.DBVideo, &video)
		if err != nil {
			log.Printf("failed with PkgVideoByVideoid: %v", err)
			return nil, nil
		}
		videos = append(videos, video)
	}
	return videos, nil
}

// 将单个DBVideo组装成FavoriteVideo
func PkgVideoByVideoid(userId int64, dbvideo data.DBVideo, video *data.FavoriteVideo) error {
	var err error
	video.Id = dbvideo.Id
	video.PlayUrl = dbvideo.PlayUrl
	video.CoverUrl = dbvideo.CoverUrl
	video.Title = dbvideo.Title
	video.Author.Id = dbvideo.AuthorId
	video.Author.Name, err = dao.SelectNameByUserId(userId)
	if err != nil {
		log.Printf("%v", err)
	}
	// video.Author.FollowCount, err = dao.CountFollow(userId)
	// if err != nil {
	// 	log.Printf("%v", err)
	// }
	// video.Author.FollowerCount, err = dao.CountFollower(userId)
	// if err != nil {
	// 	log.Printf("%v", err)
	// }
	video.Author.IsFollow, err = dao.JudgeIsFollow(userId, userId)
	if err != nil {
		log.Printf("%v", err)
	}
	// Todo
	// video.Author.Avatar = ""
	// video.Author.BackgroundImage = ""
	// video.Author.Signature = ""
	// video.Author.TotalFavorited = 0
	// video.Author.WorkCount = 0
	// video.Author.FavoriteCount = 0

	video.FavoriteCount, err = dao.CountFavorite(dbvideo.Id)
	if err != nil {
		log.Printf("failed with dao.CountFavorite(video.Id): %v", err)
	}
	video.CommentCount, err = dao.CountComment(dbvideo.Id)
	if err != nil {
		log.Printf("failed with dao.CountComment(video.Id): %v", err)
	}
	video.IsFavorite = true
	// video.IsFavorite, err = dao.JudgeIsFavorite(dbvideo.Id, userId)
	// if err != nil {
	// 	log.Printf("failed with dao.CountIsFavorite(video.Id, userID): %v", err)
	// }
	return err
}
