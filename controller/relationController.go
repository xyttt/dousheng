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

func RelationAction(c *gin.Context) {
	//获取传递的参数
	token := c.Query("token")
	to_user_id := c.Query("to_user_id")
	action_type := c.Query("action_type")
	uid, err2 := strconv.ParseInt(c.GetString("userId"), 10, 64)
	tid, err := strconv.Atoi(to_user_id)
	act, err1 := strconv.Atoi(action_type)
	if err != nil || err1 != nil || err2 != nil {
		c.JSON(http.StatusOK, data.Response{StatusCode: -1, StatusMsg: err.Error()})
		return
	}
	if len(token) == 0 || uid < 0 || act < 1 || act > 2 {
		c.JSON(http.StatusOK, data.Response{StatusCode: -1, StatusMsg: "invalid UserId or Token!"})
		return
	}
	if int64(tid) == uid {
		c.JSON(http.StatusOK, data.Response{StatusCode: 1, StatusMsg: "不能关注自己"})
		return
	}

	Req := data.DouyinRelationActionRequest{
		ToUserId:   int64(tid),
		ActionType: int32(act),
		UserId:     uid,
	}
	err = service.RelationAction(&Req)
	if err != nil {
		c.JSON(http.StatusOK, data.Response{StatusCode: -1, StatusMsg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, data.Response{
		StatusCode: 0,
		StatusMsg:  "OK",
	})
}

func FollowList(c *gin.Context) {
	curId, err1 := strconv.ParseInt(strings.TrimSpace(c.GetString("userId")), 10, 64)
	if err1 != nil {
		log.Println(err1.Error())
		c.JSON(http.StatusOK, data.Response{StatusCode: -1, StatusMsg: err1.Error()})
		return
	}
	uid, err := strconv.Atoi(c.Query("user_id")) //获取uid，当前要查看的用户
	if err != nil {
		c.JSON(http.StatusOK, data.Response{StatusCode: -1, StatusMsg: err.Error()})
		return
	}
	Userid := int64(uid)
	Token := c.Query("token")
	if len(Token) == 0 || Userid < 0 {
		c.JSON(http.StatusOK, data.Response{StatusCode: -1, StatusMsg: "invalid UserId or Token!"})
		return
	}

	resp, err := service.FollowingList(Userid, curId)

	if err != nil {
		c.JSON(http.StatusOK, data.Response{StatusCode: -1, StatusMsg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, data.DouyinRelationFollowListResponse{
		Response: data.Response{
			StatusCode: 0,
			StatusMsg:  "获取关注列表成功",
		},
		UserList: resp,
	})
	return
}

func FollowerList(c *gin.Context) {
	curId, err1 := strconv.ParseInt(strings.TrimSpace(c.GetString("userId")), 10, 64) //获取token.id-->对应当前登录用户
	if err1 != nil {
		log.Println(err1.Error())
		c.JSON(http.StatusOK, data.Response{StatusCode: -1, StatusMsg: err1.Error()})
		return
	}
	uid, err := strconv.Atoi(c.Query("user_id")) //获取uid，当前要查看的用户
	if err != nil {
		c.JSON(http.StatusOK, data.Response{StatusCode: -1, StatusMsg: err.Error()})
		return
	}
	Userid := int64(uid)

	if Userid < 0 {
		c.JSON(http.StatusOK, data.Response{StatusCode: -1, StatusMsg: "invalid UserId!"})
		return
	}

	resp, err := service.FollowerList(Userid, curId)

	if err != nil {
		c.JSON(http.StatusOK, data.Response{StatusCode: -1, StatusMsg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, data.DouyinRelationFollowListResponse{
		Response: data.Response{
			StatusCode: 0,
			StatusMsg:  "获取粉丝列表成功",
		},
		UserList: resp,
	})
	return
}

func FriendList(c *gin.Context) {
	uid, err := strconv.Atoi(c.Query("user_id"))
	if err != nil {
		c.JSON(http.StatusOK, data.Response{StatusCode: -1, StatusMsg: err.Error()})
		return
	}
	Userid := int64(uid)
	Token := c.Query("token")
	if len(Token) == 0 || Userid < 0 {
		c.JSON(http.StatusOK, data.Response{StatusCode: -1, StatusMsg: "invalid UserId or Token!"})
		return
	}
	//TODO:JWT？

	resp, err := service.FriendList(&data.DouyinRelationFollowListRequest{
		UserId: Userid, //ToUserID
		Token:  Token,
	}, Userid) //第二个参数实际上是JWT后的Userid

	if err != nil {
		c.JSON(http.StatusOK, data.Response{StatusCode: -1, StatusMsg: err.Error()})
		return
	}
	c.JSON(http.StatusOK, data.DouyinRelationFollowListResponse{
		Response: data.Response{
			StatusCode: 0,
			StatusMsg:  "获取好友列表成功",
		},
		UserList: resp,
	})
	return
}
