package data

import "time"

type User struct {
	Id            int64  `json:"id"`
	Name          string `json:"name"`
	FollowCount   int64  `json:"follow_count"`
	FollowerCount int64  `json:"follower_count"`
	IsFollow      bool   `json:"is_follow"`
}

func (User) TableName() string {
	return "users"
}

type Video struct {
	DBVideo
	Author        User  `json:"author"`
	FavoriteCount int64 `json:"favorite_count"`
	CommentCount  int64 `json:"comment_count"`
	IsFavorite    bool  `json:"is_favorite"`
}

type DBVideo struct {
	Id          int64 `json:"id"`
	AuthorId    int64
	PlayUrl     string `json:"play_url"`
	CoverUrl    string `json:"cover_url"`
	PublishTime time.Time
	Title       string `json:"title"`
}

func (DBVideo) TableName() string {
	return "videos"
}
