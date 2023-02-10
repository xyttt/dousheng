package main

import (
	"dousheng/dao"
	"dousheng/middleware/MinIO"
	"github.com/gin-gonic/gin"
)

func main() {
	dao.Init()
	dao.CreateTables()

	MinIO.Init()

	//go service.RunMessageServer() //这是什么意思？？

	r := gin.Default()

	initRouter(r)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
