package data

type Relation struct {
	ID       uint  `gorm:"primarykey"`
	UserID   int64 `gorm:"column:user_id"`
	ToUserID int64 `gorm:"column:to_user_id"`
}

func (Relation) TableName() string {
	return "follows"
}

// User Gorm Data structures
type UserRaw struct {
	ID            uint   `gorm:"primarykey"`
	Name          string `gorm:"column:name;index:idx_username,unique;type:varchar(32);not null"`
	Password      string `gorm:"column:password;type:varchar(32);not null"`
	FollowCount   int64  `gorm:"column:follow_count;default:0"`
	FollowerCount int64  `gorm:"column:follower_count;default:0"`
	Isfollow      int64
}

func (UserRaw) TableName() string {
	return "users"
}

type UserId struct {
	ID uint `gorm:"primarykey"`
}

func (UserId) TableName() string {
	return "users"
}

type DouyinRelationActionRequest struct {
	UserId     int64  `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`             // 用户id
	Token      string `protobuf:"bytes,2,opt,name=token,proto3" json:"token,omitempty"`                              // 用户鉴权token
	ToUserId   int64  `protobuf:"varint,3,opt,name=to_user_id,json=toUserId,proto3" json:"to_user_id,omitempty"`     // 对方用户id
	ActionType int32  `protobuf:"varint,4,opt,name=action_type,json=actionType,proto3" json:"action_type,omitempty"` // 1-关注，2-取消关注
}

type DouyinRelationFollowListRequest struct {
	UserId int64  `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"` // 用户id
	Token  string `protobuf:"bytes,2,opt,name=token,proto3" json:"token,omitempty"`                  // 用户鉴权token
}

type DouyinRelationFollowListResponse struct {
	Response
	UserList []*User `protobuf:"bytes,3,rep,name=user_list,json=userList,proto3" json:"user_list,omitempty"` // 用户列表
}
