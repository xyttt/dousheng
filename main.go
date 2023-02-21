package main

import (
	"dousheng/dao"
	"dousheng/middleware/MinIO"
	tool "dousheng/middleware/rabbitMQ"
	"dousheng/middleware/redis"
	"github.com/gin-gonic/gin"
)

func main() {
	dao.Init()
	dao.CreateTables()

	MinIO.Init()
	tool.InitMq()
	redis.InitRedis()
	//go service.RunMessageServer() //这是什么意思？？

	r := gin.Default()

	initRouter(r)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
