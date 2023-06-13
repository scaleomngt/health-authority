package main

import (
	"github.com/gin-gonic/gin"
	"health-info-server/service"
	"log"
)

func main() {
	log.Println("Starting")
	r := gin.Default()
	r.POST("/submitData", service.SubmitData)
	r.POST("/calcData", service.CalcData)
	r.POST("/getData", service.GetData)
	r.GET("/createAccount", service.CreateAccount)
	r.GET("/getAccounts", service.GetAccounts)
	r.GET("/initRedisId/:id", service.InitRedisId)
	r.GET("/test/:id", service.Test)
	r.Run(":7777") // 监听并在 0.0.0.0:8080 上启动服务
}
