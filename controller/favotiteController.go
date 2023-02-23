package controller

import (
	// "net/http"

	"dousheng/data"
	"dousheng/service"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type FavoriteActionResponse struct {
	StatusCode int32  `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  string `json:"status_msg"`  // 返回状态描述
}

type FavoriteListResponse struct {
	StatusCode int32                `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  string               `json:"status_msg"`  // 返回状态描述
	VideoList  []data.FavoriteVideo `json:"video_list"`  // 用户点赞视频列表
}

// FavoriteAction no practical effect, just check if token is valid
func FavoriteAction(c *gin.Context) {
	strUserId := c.GetString("user_id")
	userId, _ := strconv.ParseInt(strUserId, 10, 64)
	log.Println("userId", userId)
	strVideoId := c.Query("video_id")
	videoId, _ := strconv.ParseInt(strVideoId, 10, 64)
	strActionType := c.Query("action_type")
	actionType, _ := strconv.ParseInt(strActionType, 10, 64)

	favorite := new(service.FavoriteServiceImpl)
	//获取点赞或者取消赞操作的错误信息
	err := favorite.FavoriteAction(userId, videoId, int32(actionType))
	if err == nil {
		log.Printf("方法favorite.FavouriteAction(userid, videoId, int32(actiontype) 成功")
		c.JSON(http.StatusOK, FavoriteActionResponse{
			StatusCode: 0,
			StatusMsg:  "favourite action success",
		})
	} else {
		log.Printf("方法favorite.FavouriteAction(userid, videoId, int32(actiontype) 失败：%v", err)
		c.JSON(http.StatusOK, FavoriteActionResponse{
			StatusCode: 1,
			StatusMsg:  "favourite action fail",
		})
	}
}

// FavoriteList all users have same favorite video list
func FavoriteList(c *gin.Context) {
	strUserId := c.Query("user_id") // 传来的参数
	token := c.Query("token")
	strCurId := c.GetString("userId") // 从token中解析的参数
	userId, _ := strconv.ParseInt(strUserId, 10, 64)
	curId, _ := strconv.ParseInt(strCurId, 10, 64)
	log.Println("userId", userId)
	log.Println("curId(token)", curId)
	var videos []data.FavoriteVideo
	var err error
	favorite := new(service.FavoriteServiceImpl)
	videos, err = favorite.FavoriteList(userId, token)
	if err != nil {
		log.Printf("方法favorite.FavouriteList(userid, token) 失败")
		c.JSON(http.StatusOK, FavoriteListResponse{
			StatusCode: 1,
			StatusMsg:  "方法favorite.FavouriteList(userid, token) 失败 ",
			VideoList:  videos,
		})
	} else {
		c.JSON(http.StatusOK, FavoriteListResponse{
			StatusCode: 0,
			StatusMsg:  "方法favorite.FavouriteList(userid, token) 成功 ",
			VideoList:  videos,
		})
	}
}
