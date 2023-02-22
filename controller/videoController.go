package controller

import (
	"dousheng/data"

	"dousheng/service"

	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"strconv"
	"time"
)

type FeedResponse struct {
	data.Response
	NextTime  int64        `json:"next_time,omitempty"`
	VideoList []data.Video `json:"video_list,omitempty"` //默认空
}
type VideoListResponse struct {
	data.Response
	VideoList []data.Video `json:"video_list"`
}

// Feed:"douyin/feed/接口"
func Feed(c *gin.Context) {
	strLatestTime := c.Query("latest_time")
	log.Print("Received latest_time: " + strLatestTime)
	var latestTime time.Time
	if strLatestTime != "" {
		int64Time, _ := strconv.ParseInt(strLatestTime, 10, 64)
		latestTime = time.Unix(int64Time/1000, 0)
	} else {
		latestTime = time.Now()
	}
	log.Printf("latestTime UTS:%v", latestTime)

	strUserId := c.GetString("userId")
	userId, _ := strconv.ParseInt(strUserId, 10, 64)

	var videoService service.VideoServiceImpl

	videoList, nextTime, err := videoService.Feed(latestTime, userId)
	if err != nil {
		log.Printf("failed with videoService.Feed(latestTime, userID) : %v", err)
		c.JSON(http.StatusOK, FeedResponse{
			Response: data.Response{StatusCode: 1, StatusMsg: err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, FeedResponse{
		Response:  data.Response{StatusCode: 0, StatusMsg: "succeed"},
		NextTime:  nextTime,
		VideoList: videoList,
	})
}

func Publish(c *gin.Context) {

	videoData, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusOK, data.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	videoTitle := c.PostForm("title")
	log.Printf("title : %v", videoTitle)

	strUserId := c.GetString("userId")
	userId, _ := strconv.ParseInt(strUserId, 10, 64)
	log.Printf("user_id : %v", userId)

	var videoService service.VideoServiceImpl
	err = videoService.Publish(videoData, userId, videoTitle)
	if err != nil {
		c.JSON(http.StatusOK, data.Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, data.Response{
		StatusCode: 0,
		StatusMsg:  videoTitle + " uploaded successfully",
	})
}

func PublishList(c *gin.Context) {
	strUserId := c.GetString("userId")
	userId, _ := strconv.ParseInt(strUserId, 10, 64)
	log.Printf("user_id : %v", userId)

	var videoService service.VideoServiceImpl
	publishList, err := videoService.PubList(userId)
	if err != nil {
		log.Printf("failed with videoService.PubList(userId) : %v", err)
		c.JSON(http.StatusOK, FeedResponse{
			Response: data.Response{StatusCode: 1, StatusMsg: err.Error()},
		})
		return
	}
	c.JSON(http.StatusOK, VideoListResponse{
		Response: data.Response{
			StatusCode: 0, StatusMsg: "succeed",
		},
		VideoList: publishList,
	})
}
