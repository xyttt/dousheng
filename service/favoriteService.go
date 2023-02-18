package service

import (
	"dousheng/data"
)

type FavoriteService interface {
	FavoriteAction(user_id string, vidio_id int64, action_type int32) error
	FavoriteList(user_id int64, token string) ([]data.Video, error)
}
