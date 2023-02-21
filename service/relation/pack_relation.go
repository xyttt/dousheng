package service

import (
	"dousheng/dao"
	"dousheng/data"
	"errors"
	"gorm.io/gorm"
	"log"
)

//func GetUserInfo(fromId int64, ToUserId int64) (*data.User, error){
//	user2, err := dao.GetUserByID(ToUserId)
//	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
//		return nil, err
//	}
//	if user2 == nil {
//		return nil, err
//	}
//	user2.IsFollow = false
//	relation, err := dao.GetRelation(fromId, user2.Id)
//	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
//		return nil, err
//	}
//
//	if relation != nil {
//		user2.IsFollow = true
//	}
//	return user2, nil
//}
func PackFollowList(vs []*data.Relation, fromID int64) ([]*data.User, error) {
	users := make([]*data.User, 0)
	for _, v := range vs {
		//user2, err:= GetUserInfo(fromID, int64(v.ToUserID))
		//if err != nil {
		//	return nil, err
		//}
		user2, err := dao.GetUserByID(int64(v.ToUserID))
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		if user2 == nil {
			return nil, err
		}
		user2.IsFollow = false
		relation, err := dao.GetRelation(fromID, user2.Id)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}

		if relation != nil {
			user2.IsFollow = true
		}
		users = append(users, user2)
	}
	log.Println("从Redis中查询到所有关注者。")
	return users, nil
}

//func PackFollowerList(ctx context.Context, vs []*data.Relation, fromID int64) ([]*data.User, error) {
//	users := make([]*data.User, 0)
//	for _, v := range vs {
//		user2, err := dao.GetUserByID(ctx, int64(v.UserID))
//		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
//			return nil, err
//		}
//		users = append(users, user2)
//	}
//
//	return users, nil
//}

func PackFollowListv2(currentId int64, users []*data.UserRaw, relationMap map[int64]data.Empty) []*data.User {
	userList := make([]*data.User, 0)
	for _, user := range users {
		var isFollow bool = false

		if currentId != -1 {
			_, ok := relationMap[int64(user.ID)]
			if ok {
				isFollow = true
			}
		}
		userList = append(userList, &data.User{
			Id:            int64(user.ID),
			Name:          user.Name,
			FollowCount:   user.FollowCount,
			FollowerCount: user.FollowerCount,
			IsFollow:      isFollow,
		})
	}
	return userList
}
