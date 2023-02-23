package dao

import (
	"fmt"
	"testing"
)

func TestCommentCount(t *testing.T) {
	Init()
	count,_ := CommentCount()
	fmt.Printf("%v",count)
}

func TestInsertInsertComment(t *testing.T) {
	Init()
	CreateTables()

	user_id := int64(1)
	video_id := int64(10030)
	commentText := string("abc")

	InsertComment(user_id, video_id, commentText)
}

func  TestDeleteteComment(t *testing.T) {
	Init()
	
	err := DeleteteComment(1)
	fmt.Printf("%v", err)
}

func  TestGetCommentList(t *testing.T) {
	Init()
	list,_ := GetCommentList(10000)

	for _, com := range list {
		fmt.Printf("%v", com)
	}
	
}