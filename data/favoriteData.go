package data

type FavoriteVideo struct {
	Id            int64             `json:"id"`
	Author        FavoriteVideoUser `json:"author"`
	PlayUrl       string            `json:"play_url"`
	CoverUrl      string            `json:"cover_url"`
	FavoriteCount int64             `json:"favorite_count"`
	CommentCount  int64             `json:"comment_count"`
	IsFavorite    bool              `json:"is_favorite"`
	Title         string            `json:"title"`
}
type FavoriteVideoUser struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
	// FollowCount   int64  `json:"follow_count"`
	// FollowerCount int64  `json:"follower_count"`
	IsFollow bool `json:"is_follow"`
	// Avatar          string `json:"avatar"`
	// BackgroundImage string `json:"background_image"`
	// Signature       string `json:"signature"`
	// TotalFavorited  int64  `json:"total_favorited"`
	// WorkCount       int64  `json:"work_count"`
	// FavoriteCount   int64  `json:"favorite_count"`
}
