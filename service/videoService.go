package service

import (
	"dousheng/data"
	"mime/multipart"
	"time"
)

type VideoServiceImpl struct {
}

type VideoServicer interface {
	// 输入时间戳latestTime和userID， 返回设定长度videoList（包括用户信息）和最早发布时间作为nextTime。
	Feed(latestTime time.Time, userId int64) ([]data.Video, int64, error)

	PubList(userId int64) ([]data.Video, error)

	Publish(videoData *multipart.FileHeader, userId int64, videoTitle string) error
}
