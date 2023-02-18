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

func (favorite *FavoriteServiceImpl) FavoriteList(user_id int64, token string) ([]data.Video, error) {
	favoriteList, err1 := dao.GetFavoriteList(user_id)
	if err1 != nil {
		log.Printf("方法:GetFavouriteList query key失败: %v", err1)
		return nil, err1
	}
	favoriteVideoList, err2 := dao.GetVideoListByVideoIds(favoriteList)
	if err2 != nil {
		log.Printf("方法:GetVideoListByVideoIds query key失败: %v", err2)
		return nil, err2
	}
	return favoriteVideoList, nil
}
