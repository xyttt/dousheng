package dao

import (
	"dousheng/config"
	"dousheng/data"
	"dousheng/middleware/redis"
	"errors"
	"github.com/u2takey/go-utils/klog"
	"gorm.io/gorm"
	"log"
)

// GetRelation get relation info
func GetRelation(uid int64, tid int64) (*data.Relation, error) {
	relation := new(data.Relation)

	if err := DB.First(&relation, "user_id = ? and to_user_id = ?", uid, tid).Error; err != nil {
		return nil, err
	}
	return relation, nil
}

// 优化：
func QueryRelationByIds(currentId int64, userIds []int64) (map[int64]data.Empty, error) {
	var relations []*data.Relation
	err := DB.Where("user_id = ? AND to_user_id IN ?", currentId, userIds).Find(&relations).Error
	if err != nil {
		klog.Error("query relation by ids " + err.Error())
		return nil, err
	}
	relationMap := make(map[int64]data.Empty)
	for _, relation := range relations {
		relationMap[relation.ToUserID] = data.Empty{}
	}
	return relationMap, nil
}

func NewRelation(uid int64, tid int64) error {

	err := DB.Transaction(func(tx *gorm.DB) error {
		// 在事务中执行一些 db 操作（从这里开始，您应该使用 'tx' 而不是 'db'）
		// 1. 新增关注数据
		err := tx.Create(&data.Relation{UserID: uid, ToUserID: tid}).Error
		if err != nil {
			return err
		}

		// 2.改变 user 表中的 following count
		res := tx.Model(new(users)).Where("ID = ?", uid).Update("following_count", gorm.Expr("following_count + ?", 1))
		if res.Error != nil {
			return res.Error
		}

		if res.RowsAffected != 1 {
			return errors.New("Database error")
		}

		// 3.改变 user 表中的 follower count
		res = tx.Model(new(users)).Where("ID = ?", tid).Update("follower_count", gorm.Expr("follower_count + ?", 1))
		if res.Error != nil {
			return res.Error
		}

		if res.RowsAffected != 1 {
			return errors.New("Database error")
		}

		return nil
	})
	return err
}

// DisRelation deletes a relation from the database.
func DisRelation(uid int64, tid int64) error {
	err := DB.Transaction(func(tx *gorm.DB) error {
		// 在事务中执行一些 db 操作（从这里开始，您应该使用 'tx' 而不是 'db'）
		relation := new(data.Relation)
		if err := tx.Where("user_id = ? AND to_user_id=?", uid, tid).First(&relation).Error; err != nil {
			return err
		}

		// 1. 删除关注数据
		err := tx.Unscoped().Delete(&relation).Error
		if err != nil {
			return err
		}
		// 2.改变 user 表中的 following count
		res := tx.Model(new(users)).Where("ID = ?", uid).Update("following_count", gorm.Expr("following_count - ?", 1))
		if res.Error != nil {
			return res.Error
		}

		if res.RowsAffected != 1 {
			return errors.New("Database error")
		}

		// 3.改变 user 表中的 follower count
		res = tx.Model(new(users)).Where("ID = ?", tid).Update("follower_count", gorm.Expr("follower_count - ?", 1))
		if res.Error != nil {
			return res.Error
		}

		if res.RowsAffected != 1 {
			return errors.New("Database error")
		}

		return nil
	})
	return err
}

// FollowingList returns the Following List.
func FollowingList(uid int64) ([]*data.Relation, error) {
	var RelationList []*data.Relation
	err := DB.Debug().Where("user_id = ?", uid).Find(&RelationList).Error
	if err != nil {
		log.Println("query follow by id fail: " + err.Error())
		return nil, err
	}
	return RelationList, nil
}

// FollowerList returns the Follower List.
func FollowerList(tid int64) ([]*data.Relation, error) {
	var RelationList []*data.Relation
	//判断，如果是大V：抖音的实现就是只能看前面几十个粉丝，后面的粉丝就无法查看。
	//限制最大查询数量，如果粉丝很多,且是大V，只返回前100个粉丝，后续可根据offset一起使用
	//正常用户就正常查询即可。
	_, ok := redis.StarUsers[tid]
	if ok { //是大V，增加一个查询粉丝数的限制
		err := DB.Debug().Where("to_user_id = ?", tid).Limit(config.MaxQueryNumber).Find(&RelationList).Error
		if err != nil {
			log.Println("query StarUser fans by id fail: " + err.Error())
			return nil, err
		}
	} else {
		err := DB.Debug().Where("to_user_id = ?", tid).Find(&RelationList).Error
		if err != nil {
			log.Println("query fans by id fail: " + err.Error())
			return nil, err
		}
	}

	return RelationList, nil
}

func FriendsList(uid int64) ([]*data.Relation, error) {
	var RelationList []*data.Relation
	err := DB.Debug().Table("(?) as u",
		DB.Table("follows f1").Select("f1.user_id, f1.to_user_id").
			Joins("join follows f2 on f1.user_id = f2.to_user_id AND f1.to_user_id = f2.user_id")).
		Where("u.user_id = ?", uid).Find(&RelationList).Error
	if err != nil {
		return nil, err
	}
	return RelationList, nil
}

// select * from (select a.user_id ,a.to_user_id from follows as a inner join follows as b on a.user_id = b.to_user_id  and a.to_user_id = b.user_id) u where u.user_id = 1;
// select a.* from (select a.user_id from follower as a inner join follower as b on a.follower_id = '1' and b.follower_id = '1' ) a group by a.user_id;
// GetUserByID
func GetUserByID(userID int64) (*data.User, error) {
	res := new(data.User)
	if err := DB.First(&res, userID).Error; err != nil {
		return nil, err
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		log.Println("find no User Results!")
		return nil, err
	}
	return res, nil
}

//INSERT INTO `follower` (`id`, `user_id`, `follower_id`) VALUES  ('10', '4', '1');
//SELECT f1.*
//FROM follower f1
//JOIN follower f2 ON f1.user_id = f2.follower_id AND f1.follower_id = f2.user_id
//WHERE f1.user_id = 1

// 优化：一次获取多个Users信息
// 根据用户id获取用户信息
func QueryUserByIds(userIds []int64) ([]*data.UserRaw, error) {
	var users []*data.UserRaw
	err := DB.Where("id in (?)", userIds).Find(&users).Error
	if err != nil {
		klog.Error("query user by ids fail " + err.Error())
		return nil, err
	}
	return users, nil
}

func QueryFollow(userId int, toUserId int) (count int64, err error) {
	err = DB.Where("user_id = ? and to_user_id = ?", userId, toUserId).Find(&data.Relation{}).Count(&count).Error
	return
}
