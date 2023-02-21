package controller

import (
	"dousheng/data"
	service "dousheng/service/relation"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// 聊天记录
func MessageChat(c *gin.Context) {
	curId, err1 := strconv.ParseInt(strings.TrimSpace(c.GetString("userId")), 10, 64)
	if err1 != nil {
		log.Println(err1.Error())
		c.JSON(http.StatusOK, data.Response{StatusCode: -1, StatusMsg: err1.Error()})
		return
	}
	uid, err := strconv.Atoi(c.Query("to_user_id"))
	lastTime, err1 := strconv.Atoi(c.Query("pre_msg_time"))
	if err != nil || err1 != nil {
		c.JSON(http.StatusOK, data.Response{StatusCode: -1, StatusMsg: err.Error() + err1.Error()})
		return
	}
	ToUserid := int64(uid)
	Token := c.Query("token")
	if len(Token) == 0 || ToUserid < 0 {
		c.JSON(http.StatusOK, data.Response{StatusCode: -1, StatusMsg: "invalid UserId or Token!"})
		return
	}

	message, err := service.HistoryMessage(&data.DouyinMessageHistoryRequest{
		UserId:   curId, // TODO:修改为鉴权
		ToUserId: ToUserid,
		LastTime: int64(lastTime),
	})

	if err != nil {
		c.JSON(http.StatusOK, data.Response{StatusCode: -1, StatusMsg: err.Error()})
		return
	}

	c.JSON(http.StatusOK, data.DouyinMessageHistoryResponse{
		Response: data.Response{
			StatusCode: 0,
			StatusMsg:  "获取聊天记录成功",
		},
		Messages: message,
	})

}

// 发送消息
func MessageAction(c *gin.Context) {
	uid, err := strconv.Atoi(c.Query("user_id"))
	action, err := strconv.Atoi(c.Query("action_type"))
	content := c.Query("content")
	if err != nil {
		c.JSON(http.StatusOK, data.Response{StatusCode: -1, StatusMsg: err.Error()})
		return
	}
	Userid := int64(uid)
	Token := c.Query("token")
	if len(Token) == 0 || Userid < 0 || action != 1 {
		c.JSON(http.StatusOK, data.Response{StatusCode: -1, StatusMsg: "invalid UserId or Token!"})
		return
	}
	//TODO:JWT

	//调service层函数
	err = service.SendMessage(&data.DouyinMessageActionRequest{
		Content:  content,
		UserId:   Userid, // TODO:修改为鉴权
		ToUserId: Userid,
	})

	if err != nil {
		c.JSON(http.StatusOK, data.Response{StatusCode: -1, StatusMsg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, data.Response{StatusCode: 0, StatusMsg: "打入消息队列成功"})
	return
}
