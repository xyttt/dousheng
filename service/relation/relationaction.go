package service

import (
	"dousheng/data"
	"dousheng/middleware/rabbitMQ"
	"dousheng/middleware/redis"
	"log"
	"strconv"
	"strings"
)

func AddFollowRelation(userId int64, targetId int64) (bool, error) {
	key := strconv.Itoa(int(userId)) + strconv.Itoa(int(targetId)) + "follow"
	exist := redis.IsExistsCache(key, redis.GetRdbRelationClient())
	if exist == 1 {
		return true, nil
	}
	// 加信息打入消息队列。
	sb := strings.Builder{}
	sb.WriteString(strconv.Itoa(int(userId)))
	sb.WriteString(" ")
	sb.WriteString(strconv.Itoa(int(targetId)))
	tool.RmqFollowAdd.Publish(sb.String())
	// 记录日志
	log.Println("消息打入成功。")
	return true, nil
}
func DeleteFollowRelation(userId int64, targetId int64) (bool, error) {
	// 加信息打入消息队列。
	sb := strings.Builder{}
	sb.WriteString(strconv.Itoa(int(userId)))
	sb.WriteString(" ")
	sb.WriteString(strconv.Itoa(int(targetId)))
	tool.RmqFollowDel.Publish(sb.String())
	// 记录日志
	log.Println("消息打入成功。")
	// 更新redis信息。
	return updateRedisWithDel(userId, targetId)
}
func RelationAction(req *data.DouyinRelationActionRequest) error {
	// 1-关注
	if req.ActionType == 1 {
		go AddFollowRelation(req.UserId, req.ToUserId)
	}
	// 2-取消关注
	if req.ActionType == 2 {
		go DeleteFollowRelation(req.UserId, req.ToUserId)
	}
	return nil
}

func updateRedisWithAdd(userId int64, targetId int64) (bool, error) {
	return true, nil
}

func updateRedisWithDel(userId int64, targetId int64) (bool, error) {
	return true, nil
}
