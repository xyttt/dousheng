package dao

import (
	"dousheng/config"
	"dousheng/data"
)

// type Favorite struct {
// 	ID         int64 `gorm:"primarykey"` // auto increment primary key
// 	UserId     int64 // user_id who click favorite
// 	VideoId    int64 // video_id which is clicked favorite
// 	ActionType int32 // favorite or not, 1 is for like, 2 is for unlike
// }

func (f Favorite) TableName() string {
	return "favorites"
}

func GetFavorite(user_id int64, video_id int64) (Favorite, error) {
	var favorite Favorite
	result := DB.Model(Favorite{}).Where(&Favorite{UserId: user_id, VideoId: video_id}).Find(&favorite)
	// 处理错误
	if result.Error != nil {
		return favorite, result.Error
	}
	return favorite, nil
}

func InsertFavorite(user_id int64, video_id int64) error {
	favorite := &Favorite{UserId: user_id, VideoId: video_id, ActionType: 1}
	err := DB.Model(Favorite{}).Create(favorite).Error
	return err
}
func UpdateFavorite(user_id int64, video_id int64, action_type int32) error {
	result := DB.Model(Favorite{}).Where(map[string]interface{}{"user_id": user_id, "video_id": video_id}).
		Update("action_type", action_type)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
func GetFavoriteList(user_id int64) ([]int64, error) {
	var favoriteList []int64
	result := DB.Model(Favorite{}).Where(map[string]interface{}{"user_id": user_id, "action_type": 1}).Pluck("video_id", &favoriteList)
	if result.Error != nil {
		return favoriteList, result.Error
	}
	return favoriteList, nil
}

// 返回的是DBVideo List
func GetVideoListByVideoIds(favoriteList []int64) ([]data.Video, error) {
	var videos []data.Video
	db := GetDB()
	res := db.Table("videos").Where("id in (?)", favoriteList).Order("publish_time desc").Limit(config.VideoNum).Find(&videos)
	return videos, res.Error
}
