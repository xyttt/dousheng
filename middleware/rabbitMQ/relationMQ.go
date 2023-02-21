package tool

import (
	"dousheng/dao"
	"dousheng/middleware/redis"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"strconv"
	"strings"
	"time"
)

type FollowMQ struct {
	RabbitMQ
	channel   *amqp.Channel
	queueName string
	exchange  string
	key       string
}

// NewFollowRabbitMQ 获取followMQ的对应队列。
func NewFollowRabbitMQ(queueName string) *FollowMQ {
	followMQ := &FollowMQ{
		RabbitMQ:  *Rmq,
		queueName: queueName,
	}

	cha, err := followMQ.conn.Channel()
	followMQ.channel = cha
	Rmq.failOnErr(err, "获取通道失败")
	return followMQ
}

// 关闭mq通道和mq的连接。
func (f *FollowMQ) destroy() {
	f.channel.Close()
}

// Publish follow关系的发布配置。
func (f *FollowMQ) Publish(message string) {

	_, err := f.channel.QueueDeclare(
		f.queueName,
		//是否持久化
		false,
		//是否为自动删除
		false,
		//是否具有排他性
		false,
		//是否阻塞
		false,
		//额外属性
		nil,
	)
	if err != nil {
		panic(err)
	}

	f.channel.Publish(
		f.exchange,
		f.queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})

}

// Consumer follow关系的消费逻辑。
func (f *FollowMQ) Consumer() {

	_, err := f.channel.QueueDeclare(f.queueName, false, false, false, false, nil)

	if err != nil {
		panic(err)
	}

	//2、接收消息
	msgs, err := f.channel.Consume(
		f.queueName,
		//用来区分多个消费者
		"",
		//是否自动应答
		true,
		//是否具有排他性
		false,
		//如果设置为true，表示不能将同一个connection中发送的消息传递给这个connection中的消费者
		false,
		//消息队列是否阻塞
		false,
		nil,
	)
	if err != nil {
		panic(err)
	}

	forever := make(chan bool)
	switch f.queueName {
	case "follow_add":
		go f.consumerFollowAdd(msgs)
	case "follow_del":
		go f.consumerFollowDel(msgs)

	}

	log.Printf("[*] Waiting for messagees,To exit press CTRL+C")

	<-forever

}

// 关系添加的消费方式。
func (f *FollowMQ) consumerFollowAdd(msgs <-chan amqp.Delivery) {
	for d := range msgs {
		// 参数解析。
		//缓存延时双删
		params := strings.Split(fmt.Sprintf("%s", d.Body), " ")
		userId, _ := strconv.Atoi(params[0])
		targetId, _ := strconv.Atoi(params[1])
		key := strconv.Itoa(userId) + strconv.Itoa(targetId) + "follow"
		//查询数据库是否存在关系：若存在，什么也不做。
		count, err := dao.QueryFollow(userId, targetId)
		if err != nil {
			log.Println("查询关注信息失败", err)
		}
		if count == 0 {
			log.Println("准备延时双删")
			err = redis.DelCache(fmt.Sprintf("followList%v", userId), redis.GetRdbFollowingClient())
			err = redis.DelCache(fmt.Sprintf("fanList%v", userId), redis.GetRdbFollowerClient())
			if err := dao.NewRelation((int64)(userId), (int64)(targetId)); err != nil {
				log.Println(err.Error())
			}

			go func() {
				time.Sleep(time.Millisecond) //延时双删策略。保证数据一致性。
				err = redis.DelCache(fmt.Sprintf("followList%v", userId), redis.GetRdbFollowingClient())
				err = redis.DelCache(fmt.Sprintf("fanList%v", userId), redis.GetRdbFollowerClient())
				redis.SetRedisNum(key, key, redis.GetRdbRelationClient())
			}()
		}

	}
}

// 关系删除的消费方式 ---- 缓存延时双删
func (f *FollowMQ) consumerFollowDel(msgs <-chan amqp.Delivery) {
	for d := range msgs {
		// 参数解析。
		params := strings.Split(fmt.Sprintf("%s", d.Body), " ")
		userId, _ := strconv.Atoi(params[0])
		targetId, _ := strconv.Atoi(params[1])

		key := strconv.Itoa(userId) + strconv.Itoa(targetId) + "follow"
		//查询数据库是否存在关系：若存在，什么也不做。
		count, err := dao.QueryFollow(userId, targetId)
		if err != nil {
			log.Println("查询关注信息失败", err)
		}
		if count == 1 {
			log.Println("准备延时双删")
			key1 := fmt.Sprintf("followList%v", userId)
			key2 := fmt.Sprintf("fansList%v", userId)
			err = redis.DelCache(key1, redis.GetRdbFollowingClient())
			err = redis.DelCache(key2, redis.GetRdbFollowerClient())
			if err := dao.DisRelation((int64)(userId), (int64)(targetId)); err != nil {
				log.Println(err.Error())
			}

			go func() {
				time.Sleep(time.Millisecond) //延时双删策略。保证数据一致性。
				err = redis.DelCache(key1, redis.GetRdbFollowingClient())
				err = redis.DelCache(key2, redis.GetRdbFollowerClient())
				redis.DelCache(key, redis.GetRdbRelationClient())
			}()
		}

	}
}

var RmqFollowAdd *FollowMQ
var RmqFollowDel *FollowMQ

// InitFollowRabbitMQ 初始化rabbitMQ连接。
func InitFollowRabbitMQ() {
	RmqFollowAdd = NewFollowRabbitMQ("follow_add")
	go RmqFollowAdd.Consumer()

	RmqFollowDel = NewFollowRabbitMQ("follow_del")
	go RmqFollowDel.Consumer()
}
