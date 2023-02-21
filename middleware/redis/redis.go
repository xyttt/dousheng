package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

var Ctx = context.Background()
var RdbFollowers *redis.Client
var RdbFollowing *redis.Client
var RdbFriends *redis.Client
var RdbRelations *redis.Client

func InitRedis() {
	RdbFollowers = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "tiktok",
		DB:       0, // 粉丝列表信息存入 DB0.
		//连接池容量及闲置连接数量
		PoolSize:     15, // 连接池最大socket连接数，默认为4倍CPU数， 4 * runtime.NumCPU
		MinIdleConns: 10, //在启动阶段创建指定数量的Idle连接，并长期维持idle状态的连接数不少于指定数量；。

		//超时
		DialTimeout:  5 * time.Second, //连接建立超时时间，默认5秒。
		ReadTimeout:  3 * time.Second, //读超时，默认3秒， -1表示取消读超时
		WriteTimeout: 3 * time.Second, //写超时，默认等于读超时
		PoolTimeout:  4 * time.Second, //当所有连接都处在繁忙状态时，客户端等待可用连接的最大等待时长，默认为读超时+1秒。

		//闲置连接检查包括IdleTimeout，MaxConnAge
		IdleCheckFrequency: 60 * time.Second, //闲置连接检查的周期，默认为1分钟，-1表示不做周期性检查，只在客户端获取连接时对闲置连接进行处理。
		IdleTimeout:        5 * time.Minute,  //闲置超时，默认5分钟，-1表示取消闲置超时检查
		MaxConnAge:         0 * time.Second,  //连接存活时长，从创建开始计时，超过指定时长则关闭连接，默认为0，即不关闭存活时长较长的连接

	})
	RdbFollowing = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "tiktok",
		DB:       1, // 关注列表信息信息存入 DB1.
	})
	RdbFriends = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "tiktok",
		DB:       2, // 当前用户是否关注了自己粉丝信息存入 DB1.
	})
	RdbRelations = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "tiktok",
		DB:       3, // 当前用户是否关注了自己粉丝信息存入 DB1.
	})

}

func GetRdbFriendsClient() *redis.Client {
	return RdbFriends
}

func GetRdbFollowingClient() *redis.Client {
	return RdbFollowing
}

func GetRdbFollowerClient() *redis.Client {
	return RdbFollowers
}

func GetRdbRelationClient() *redis.Client {
	return RdbRelations
}
