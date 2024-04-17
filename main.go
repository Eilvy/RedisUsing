package main

import (
	"RedisUsing/config"
	"RedisUsing/test"
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

func main() {
	router := gin.Default() //启动一个gin服务
	config.RedisConnect()   //连接到redis数据库
	//config.DB.Set()
	//fmt.Println(config.DB.Get(context.Background(), "name"))
	//未设置中间件
	router.POST("/creat/User", config.CreateUser)      //创建用户和分数
	router.GET("/revrangeUsers", config.RevrangeUsers) //排序请求

	router.POST("/publish/:channel", config.PublishHandler)    //发布订阅频道信息
	router.GET("/subscribe/:channel", config.SubscribeHandler) //接收订阅频道信息

	err := router.Run(":8080")
	if err != nil {
		fmt.Println("启动gin服务失败")
		return
	}

	//开启两个task
	key := "daseradfse"
	expireTime := 5 * time.Second
	config.Wg.Add(2)

	go test.Task1(key, expireTime, &config.Wg)
	go test.Task2(key, expireTime, &config.Wg)

	config.Wg.Wait()

}
