package dao

import (
	"fmt"
	"testing"
)

func TestInsertFavorite(t *testing.T) {
	Init()
	CreateTables()
	user_id := int64(3)
	// video_id := int64(10010)
	video_id := int64(12345)

	err := InsertFavorite(user_id, video_id)
	fmt.Printf("%v", err)
}

func TestGetFavorite(t *testing.T) {
	Init()
	// user_id := int64(1)
	user_id := int64(3)
	video_id := int64(10010)
	favorite, err := GetFavorite(user_id, video_id)
	fmt.Println(favorite.Id)
	fmt.Printf("%v", err)
}

func TestUpdateFavorite(t *testing.T) {
	Init()
	user_id := int64(1)
	video_id := int64(10010)
	action_type := int32(2)
	err := UpdateFavorite(user_id, video_id, action_type)
	fmt.Printf("%v", err)
}

func TestGetFavoriteList(t *testing.T) {
	Init()
	user_id := int64(1)
	favorite, err := GetFavoriteList(user_id)
	fmt.Println(favorite)
	fmt.Printf("%v", err)
}

func TestGetVideoListByVideoIds(t *testing.T) {
	Init()
	user_id := int64(0)
	favoriteList, err := GetFavoriteList(user_id)
	// var favoriteVideoList []data.Video
	// var e error
	favoriteVideoList, e := GetVideoListByVideoIds(favoriteList)
	fmt.Println(favoriteVideoList)
	fmt.Printf("%v", err)
	fmt.Printf("%v", e)
}
