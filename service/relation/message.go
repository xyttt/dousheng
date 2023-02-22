package service

import (
	"dousheng/dao"
	"dousheng/data"
	tool "dousheng/middleware/rabbitMQ"
	"log"
	"strconv"
	"strings"
)

func SendMessage(res *data.DouyinMessageActionRequest) error {
	sb := strings.Builder{}
	_, err := sb.WriteString(strconv.Itoa(int(res.UserId)))
	if err != nil {
		return err
	}
	sb.WriteString(" ")
	_, err1 := sb.WriteString(strconv.Itoa(int(res.ToUserId)))
	if err1 != nil {
		return err1
	}
	sb.WriteString(" ")
	_, err2 := sb.WriteString(res.Content)
	if err2 != nil {
		return err2
	}
	tool.RmqMeaasge.Publish(sb.String())
	log.Println("消息打入成功。")
	return nil
}

func HistoryMessage(res *data.DouyinMessageHistoryRequest) ([]*data.Message, error) {
	messages, err := dao.MessageHistory(res.UserId, res.ToUserId, res.LastTime)
	if err != nil {
		return nil, err
	}
	return messages, nil
}
