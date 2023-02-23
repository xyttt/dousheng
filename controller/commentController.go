package controller

import (
	"dousheng/data"
	"dousheng/service"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CommentActionResponse struct {
	StatusCode int32        `json:"status_code"`
	StatusMsg  string       `json:"status_msg"`
	Comment    data.Comment `json:"comment"`
}

type CommentListResponse struct {
	StatusCode   int32          `json:"status_code"`
	StatusMsg    string         `json:"status_msg"`
	Comment_list []data.Comment `json:"comment_list"`
}

func CommentAction(c *gin.Context) {
	// strToken := c.Query("token")
	// userToken, _ := strconv.ParseInt(strToken, 10, 64)
	// userId := userToken
	strUserId := c.GetString("user_id")
	userId, _ := strconv.ParseInt(strUserId, 10, 64)
	log.Printf("userId: %v", userId)

	strVideoId := c.Query("video_id")
	videoId, _ := strconv.ParseInt(strVideoId, 10, 64)
	log.Printf("videoId: %v", videoId)

	strActionType := c.Query("action_type")
	actionType, _ := strconv.ParseInt(strActionType, 10, 64)
	log.Printf("actionType: %v", actionType)

	commentText := ""
	if actionType == 1 {
		commentT := c.Query("comment_text")
		commentText = commentT
	}
	log.Printf("commentText: %v", commentText)

	commentId := int64(0)
	if actionType == 2 {
		strCommentId := c.Query("comment_id")
		commentI, _ := strconv.ParseInt(strCommentId, 10, 64)
		commentId = commentI
	}
	log.Printf("commentId: %v", commentId)

	// commentText := c.GetString("comment_text")
	// log.Printf("commentText: %v", commentText)

	// strCommentId := c.Query("comment_id")
	// commentId, _ := strconv.ParseInt(strCommentId, 10, 64)
	// log.Printf("commentId: %v", commentId)

	comment := new(service.CommentServiceImpl)
	com, err := comment.CommentAction(userId, videoId, int32(actionType), commentText, commentId)

	if err == nil {
		c.JSON(http.StatusOK, CommentActionResponse{
			StatusCode: 0,
			StatusMsg:  "comment action success",
			Comment:    com,
		})
	} else {
		c.JSON(http.StatusOK, CommentActionResponse{
			StatusCode: 1,
			StatusMsg:  "comment action fail",
			Comment:    com,
		})
	}
}

func CommentList(c *gin.Context) {
	c.Query("token")
	// strToken := c.Query("token")
	// userToken, _ := strconv.ParseInt(strToken, 10, 64)
	// userId := userToken

	strVideoId := c.Query("video_id")
	videoId, _ := strconv.ParseInt(strVideoId, 10, 64)

	comment := new(service.CommentServiceImpl)
	var coms []data.Comment
	var err error
	coms, err = comment.CommentList(videoId)
	if err != nil {
		c.JSON(http.StatusOK, CommentListResponse{
			StatusCode:   1,
			StatusMsg:    "comment list fail ",
			Comment_list: coms,
		})
	} else {
		c.JSON(http.StatusOK, CommentListResponse{
			StatusCode:   0,
			StatusMsg:    "comment list success ",
			Comment_list: coms,
		})
	}
}
