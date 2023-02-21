package tool

import (
	"dousheng/dao"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"strconv"
	"strings"
)

type MessageMQ struct {
	RabbitMQ
	channel   *amqp.Channel
	queueName string
	exchange  string
	key       string
}

// NewFollowRabbitMQ 获取followMQ的对应队列。
func NewMessageRabbitMQ(queueName string) *MessageMQ {
	meaasgeMQ := &MessageMQ{
		RabbitMQ:  *Rmq,
		queueName: queueName,
	}

	cha, err := meaasgeMQ.conn.Channel()
	meaasgeMQ.channel = cha
	Rmq.failOnErr(err, "获取通道失败")
	return meaasgeMQ
}

// 关闭mq通道和mq的连接。
func (f *MessageMQ) destroy() {
	f.channel.Close()
}

// Publish follow关系的发布配置。
func (f *MessageMQ) Publish(message string) {

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
func (f *MessageMQ) Consumer() {

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
	go f.messageAdd(msgs)

	log.Printf("[*] Waiting for messagees,To exit press CTRL+C")

	<-forever

}

// 消息添加的消费方式。
func (f *MessageMQ) messageAdd(msgs <-chan amqp.Delivery) {
	for d := range msgs {

		params := strings.Split(fmt.Sprintf("%s", d.Body), " ")
		userId, _ := strconv.Atoi(params[0])
		targetId, _ := strconv.Atoi(params[1])
		content := params[2]

		if err := dao.CreateMes((int64)(userId), (int64)(targetId), content); err != nil {
			log.Println(err.Error())
		}

	}
}

var RmqMeaasge *MessageMQ

// InitFollowRabbitMQ 初始化rabbitMQ连接。
func InitMessageRabbitMQ() {
	RmqMeaasge = NewMessageRabbitMQ("message_add")
	go RmqMeaasge.Consumer()
}
