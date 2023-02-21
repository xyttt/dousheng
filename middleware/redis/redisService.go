package redis

import (
	"bytes"
	"crypto/rand"
	"dousheng/data"
	"encoding/gob"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/u2takey/go-utils/klog"
	"log"
	"math"
	"math/big"
	"time"
)

func IsExistsCache(key string, conn *redis.Client) (exists int64) {
	exist, err := conn.Exists(Ctx, key).Result()
	if err != nil {
		klog.Error("查询缓存是否存在失败", err)
	}
	return exist
}

func GetFanListCache(userId string) (result []*data.User, err error) {
	conn := GetRdbFollowerClient()
	key := fmt.Sprintf("fansList%v", userId)
	rebytes, err := conn.Get(Ctx, key).Result()
	if err != nil {
		log.Printf("读取fansList%v缓存失败,err:%v", userId, err)
	}

	var fansList []data.User
	//进行gob序列化
	reader := bytes.NewReader([]byte(rebytes))
	dec := gob.NewDecoder(reader)
	err = dec.Decode(&fansList)

	//转换为指针传递。
	result = make([]*data.User, len(fansList))
	for i := 0; i < len(result); i++ {
		result[i] = &fansList[i]
	}
	return
}
func GetFollowingCache(userId string) (result []*data.User, err error) {
	conn := GetRdbFollowingClient()
	key := fmt.Sprintf("followList%v", userId)
	rebytes, err := conn.Get(Ctx, key).Result()
	if err != nil {
		log.Printf("读取followList%v缓存失败,err:%v", userId, err)
	}

	var followingList []data.User
	//进行gob序列化
	reader := bytes.NewReader([]byte(rebytes))
	dec := gob.NewDecoder(reader)
	err = dec.Decode(&followingList)

	//转换为指针传递。
	result = make([]*data.User, len(followingList))
	for i := 0; i < len(result); i++ {
		result[i] = &followingList[i]
	}
	return
}

func GetFriendsListCache(userId string) (result []*data.User, err error) {
	conn := GetRdbFriendsClient()
	key := fmt.Sprintf("friendsList%v", userId)
	rebytes, err := conn.Get(Ctx, key).Result()
	if err != nil {
		log.Printf("读取friendsList%v缓存失败,err:%v", userId, err)
	}

	var friendsList []data.User
	//进行gob序列化
	reader := bytes.NewReader([]byte(rebytes))
	dec := gob.NewDecoder(reader)
	err = dec.Decode(&friendsList)

	//转换为指针传递。
	result = make([]*data.User, len(friendsList))
	for i := 0; i < len(result); i++ {
		result[i] = &friendsList[i]
	}
	return
}

func SetNull(key string, conn *redis.Client) (err error) {
	_, err = conn.SetEX(Ctx, key, "", time.Duration(10)).Result()
	if err != nil {
		klog.Error("缓存空值到%s失败,err:%v", key, err)
	}
	return
}
func SetRedisNum(key, value string, conn *redis.Client) {
	randNum := time.Duration(rangeRand(30, 60)) * time.Minute
	times := 10*time.Hour + randNum
	_, err := conn.SetEX(Ctx, key, value, times).Result()
	if err != nil {
		klog.Error("设置缓存失败", err)
	}
}

func rangeRand(min, max int64) int64 {
	if min > max {
		panic("the min is greater than max!")
	}
	if min < 0 {
		f64Min := math.Abs(float64(min))
		i64Min := int64(f64Min)
		result, _ := rand.Int(rand.Reader, big.NewInt(max+1+i64Min))
		return result.Int64() - i64Min
	} else {
		result, _ := rand.Int(rand.Reader, big.NewInt(max-min+1))
		return min + result.Int64()
	}
}

//使用泛型转化所有的类型
func Convert[T any](s []*T) (res []T) {
	klog.Info("length:", len(s))
	res = make([]T, len(s))
	for i := 0; i < len(s); i++ {
		res[i] = *s[i]
	}
	return
}

func SetRedisCache(key string, data interface{}, conn *redis.Client) (err error) {
	log.Printf("starting SetRedis")
	var buffer bytes.Buffer
	ecoder := gob.NewEncoder(&buffer)
	err = ecoder.Encode(data)
	if err != nil {
		log.Println(err)
		return
	}
	//加上随机数，防止同时过期造成缓存雪崩

	randNum := time.Duration(rangeRand(10, 60)) * time.Minute
	keeptime := 10*time.Hour + randNum //10小时+半小时~1小时随机过期时间
	_, err = conn.SetEX(Ctx, key, buffer.Bytes(), keeptime).Result()
	if err != nil {
		log.Printf("写入%s缓存失败,err:%v", key, err)
	}
	return
}

func DelCache(key string, conn *redis.Client) (err error) {
	log.Print("del key:", key)
	_, err = conn.Del(Ctx, key).Result()
	if err != nil {
		klog.Info("删除%s缓存失败,err:%v", key, err)
	}
	return
}

//基于Redis实现简单的分布式锁
func RedisLock(key string, uuid string, conn *redis.Client) (isLock bool) {
	redisLockTimeout := 10
	lockSuccess, err := conn.SetNX(Ctx, key, uuid, time.Duration(redisLockTimeout)).Result()
	if err != nil || !lockSuccess {
		klog.Error("加锁失败", err)
		return false
	} else {
		klog.Info("get lock success")
	}
	return true
}

//https://blog.csdn.net/qq_27176027/article/details/126275178
func RedisUnlock(key string, uuid string, conn *redis.Client) (err error) {
	//version1：可能有并发安全问题。判断锁成功了，但是这时候这个锁过期了，其他的线程获得锁并加锁了， 这时候就会把其他的线程加的锁误删了
	//value, _ := conn.Get(Ctx, key).Result()
	//if value == uuid {
	//	_, err = conn.Del(Ctx, key).Result()
	//	if err != nil {
	//		klog.Error("解锁失败", err)
	//		return
	//	}
	//	klog.Info("解锁成功！")
	//}

	//version2： LUA脚本:确保原子性。
	script := "if redis.call('get', KEYS[1]) == ARGV[1] then return redis.call('del', KEYS[1]) else return 0 end"
	result, err := conn.Do(Ctx, "EVAL", script, 1, key, uuid).Bool()
	if !result {
		return errors.New("出现分布式并发释放锁错误")
	}
	return
}
