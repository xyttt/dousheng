package service

import (
	"dousheng/dao"
	"dousheng/data"
	"dousheng/middleware/redis"
	"errors"
	"fmt"
	"github.com/satori/go.uuid"
	"log"
	"strconv"
	"time"
)

func FollowingMySQLService(targetId int64, fromID int64) ([]*data.User, error) {
	//判断当前用户是否合法
	user, err := dao.QueryUserByIds([]int64{targetId})
	if err != nil {
		return nil, err
	}
	if len(user) == 0 {
		return nil, errors.New("userId not exist")
	}
	//查询有谁关注了该用户
	FollowingUser, err := dao.FollowingList(targetId)
	if err != nil {
		return nil, err
	}
	userIds := make([]int64, 0) //关注方信息数组
	for _, relation := range FollowingUser {
		userIds = append(userIds, relation.ToUserID)
	}
	//获取关注方信息
	users, err := dao.QueryUserByIds(userIds)
	if err != nil {
		return nil, err
	}

	var relationMap map[int64]data.Empty //利用空结构体 + Map实现Set，减少内存占用。
	if fromID == -1 {
		relationMap = nil
	} else {
		//获取当前用户与关注方的关注记录
		relationMap, err = dao.QueryRelationByIds(fromID, userIds)
		if err != nil {
			return nil, err
		}
	}
	FollowUserList := PackFollowListv2(fromID, users, relationMap)
	return FollowUserList, nil
}

func FollowingList(targetId, fromID int64) (FollowingUserList []*data.User, err error) {
	key := fmt.Sprintf("followList%v", targetId)
	if redis.IsExistsCache(key, redis.GetRdbFollowingClient()) == 1 {
		FollowingUserList, err = redis.GetFollowingCache(strconv.FormatInt(targetId, 10))
		if err != nil {
			log.Println("查询关注列表缓存失败", err)
		}
		return FollowingUserList, err
	} else {
		Id := uuid.NewV4().String() //生成这个锁的唯一识别码
		lockNum := "1"
		if redis.RedisLock(lockNum, Id, redis.GetRdbFollowingClient()) == true {
			//获取锁成功：查数据库、同步缓存、删锁
			FollowingUserList, err = FollowingMySQLService(targetId, fromID)
			if err != nil {
				log.Println("获取关注列表失败", err)
				return nil, err
			}
			if len(FollowingUserList) == 0 {
				go redis.SetNull(key, redis.GetRdbFollowingClient())
				//如果数据库不存在，则缓存一个10秒的空值，防止缓存穿透
			} else {
				go redis.SetRedisCache(key, redis.Convert(FollowingUserList), redis.GetRdbFollowerClient())
			}
			redis.RedisUnlock(lockNum, Id, redis.GetRdbFollowerClient())
			return FollowingUserList, nil
		} else { //获取锁失败，睡眠一段时间重新查询,可能会出现数据还没写完就去读了
			for i := 0; i < 100; i++ {
				time.Sleep(time.Millisecond * 100)
				FollowingUserList, err = redis.GetFollowingCache(strconv.FormatInt(targetId, 10))
				if err != nil {
					log.Println("查询点赞列表缓存失败", err)
				} else {
					return
				}
			}
			return nil, errors.New("读取数据超过设置最大次数限制")
		}
	}
	return
}

func FollowerMySQLService(targetId int64, fromID int64) ([]*data.User, error) {
	//判断当前用户是否合法
	user, err := dao.QueryUserByIds([]int64{targetId})
	if err != nil {
		return nil, err
	}
	if len(user) == 0 {
		return nil, errors.New("userId not exist")
	}
	//查询有谁关注了该用户
	FollowerUser, err := dao.FollowerList(targetId)
	if err != nil {
		return nil, err
	}
	userIds := make([]int64, 0) //关注方信息数组
	for _, relation := range FollowerUser {
		userIds = append(userIds, int64(relation.UserID))
	}
	//获取关注方信息
	users, err := dao.QueryUserByIds(userIds)
	if err != nil {
		return nil, err
	}

	var relationMap map[int64]data.Empty //李利用空结构体 + Map实现Set，减少内存占用。
	if fromID == -1 {
		relationMap = nil
	} else {
		//获取当前用户与关注方的关注记录
		relationMap, err = dao.QueryRelationByIds(fromID, userIds)
		if err != nil {
			return nil, err
		}
	}
	FollowUserList := PackFollowListv2(fromID, users, relationMap)
	return FollowUserList, nil
}

// FollowerList returns the Follower Lists
func FollowerList(targetId int64, fromID int64) (FollowUserList []*data.User, err error) {
	//先查询缓存:
	key := fmt.Sprintf("fansList%v", targetId)
	if redis.IsExistsCache(key, redis.GetRdbFollowerClient()) == 1 {
		FollowUserList, err = redis.GetFanListCache(strconv.FormatInt(targetId, 10))
		if err != nil {
			log.Println("查询粉丝列表缓存失败", err)
		}
		return FollowUserList, err
	} else {
		Id := uuid.NewV4().String() //生成这个锁的唯一识别码
		lockNum := "1"
		if redis.RedisLock(lockNum, Id, redis.GetRdbFollowerClient()) == true {
			//获取锁成功：查数据库、同步缓存、删锁
			FollowUserList, err = FollowerMySQLService(targetId, fromID)
			if err != nil {
				log.Println("获取关注列表失败", err)
				return nil, err
			}
			if len(FollowUserList) == 0 {
				go redis.SetNull(key, redis.GetRdbFollowerClient())
				//如果数据库不存在，则缓存一个10秒的空值，防止缓存穿透
			} else {
				go redis.SetRedisCache(key, redis.Convert(FollowUserList), redis.GetRdbFollowerClient())
			}
			redis.RedisUnlock(lockNum, Id, redis.GetRdbFollowerClient())
			return FollowUserList, nil
		} else { //获取锁失败，睡眠一段时间重新查询,
			for i := 0; i < 100; i++ {
				time.Sleep(time.Millisecond * 100)
				FollowUserList, err = redis.GetFanListCache(strconv.FormatInt(targetId, 10))
				if err != nil {
					log.Println("查询点赞列表缓存失败", err)
				} else {
					return
				}
			}
			return nil, errors.New("读取数据超过设置最大次数限制")
		}
	}
	return
}

func FriendsMySQLService(targetId int64, fromID int64) ([]*data.User, error) {
	//判断当前用户是否合法
	user, err := dao.QueryUserByIds([]int64{targetId})
	if err != nil {
		return nil, err
	}
	if len(user) == 0 {
		return nil, errors.New("userId not exist")
	}
	//查询有谁关注了该用户
	FriendsUser, err := dao.FriendsList(targetId)
	if err != nil {
		return nil, err
	}
	userIds := make([]int64, 0) //关注方信息数组
	for _, relation := range FriendsUser {
		userIds = append(userIds, relation.ToUserID)
	}
	//获取关注方信息
	users, err := dao.QueryUserByIds(userIds)
	if err != nil {
		return nil, err
	}

	var relationMap map[int64]data.Empty //李利用空结构体 + Map实现Set，减少内存占用。
	if fromID == -1 {
		relationMap = nil
	} else {
		//获取当前用户与关注方的关注记录
		relationMap, err = dao.QueryRelationByIds(fromID, userIds)
		if err != nil {
			return nil, err
		}
	}
	FollowUserList := PackFollowListv2(fromID, users, relationMap)
	return FollowUserList, nil
}

func FriendList(targetId int64, fromID int64) (FriendsUserList []*data.User, err error) {
	//先查询缓存:
	key := fmt.Sprintf("friendsList%v", targetId)
	if redis.IsExistsCache(key, redis.GetRdbFriendsClient()) == 1 {
		FriendsUserList, err = redis.GetFriendsListCache(strconv.FormatInt(targetId, 10))
		if err != nil {
			log.Println("查询粉丝列表缓存失败", err)
		}
		return FriendsUserList, err
	} else {
		Id := uuid.NewV4().String() //生成这个锁的唯一识别码
		lockNum := "1"
		if redis.RedisLock(lockNum, Id, redis.GetRdbFriendsClient()) == true {
			//获取锁成功：查数据库、同步缓存、删锁
			FriendsUserList, err = FriendsMySQLService(targetId, fromID)
			if err != nil {
				log.Println("获取关注列表失败", err)
				return nil, err
			}
			if len(FriendsUserList) == 0 {
				go redis.SetNull(key, redis.GetRdbFriendsClient())
				//如果数据库不存在，则缓存一个10秒的空值，防止缓存穿透
			} else {
				go redis.SetRedisCache(key, redis.Convert(FriendsUserList), redis.GetRdbFriendsClient())
			}
			redis.RedisUnlock(lockNum, Id, redis.GetRdbFriendsClient())
			return FriendsUserList, nil
		} else { //获取锁失败，睡眠一段时间重新查询,可能会出现数据还没写完就去读了
			for i := 0; i < 100; i++ {
				time.Sleep(time.Millisecond * 100)
				FriendsUserList, err = redis.GetFriendsListCache(strconv.FormatInt(targetId, 10))
				if err != nil {
					log.Println("查询点赞列表缓存失败", err)
				} else {
					return
				}
			}
			return nil, errors.New("读取数据超过设置最大次数限制")
		}
	}
	return

}
