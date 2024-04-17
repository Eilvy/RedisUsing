package config

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"net/http"
	"strconv"
)

func CreateUser(c *gin.Context) {
	var user User
	err := c.ShouldBind(&user)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	number, err := strconv.ParseFloat(user.Number, 64) //将user的string类型的number转化为float64类型
	if number > 200 {
		c.JSON(400, gin.H{"msg": "分数不能超过200"})
		return
	}
	numAdded, err := DB.ZAdd(context.Background(), "Users", redis.Z{Score: number, Member: user.Username}).Result()
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if numAdded == 0 {
		c.JSON(200, gin.H{"msg": "已添加该用户"})
		return
	}
	c.JSON(200, gin.H{"status": "success"})
}
func RevrangeUsers(c *gin.Context) {
	members, err := DB.ZRevRangeWithScores(context.Background(), "Users", 0, -1).Result()
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	var results []map[string]interface{}
	for _, member := range members { //从redis里面循环zset Users里面的成员和分数
		result := map[string]interface{}{ //创建一个新的map储存当前循环的成员和分数
			"member": member.Member,
			"score":  member.Score,
		}
		results = append(results, result)
	}
	c.JSON(http.StatusOK, gin.H{"members": results})
} //对用户按照分数从高到低进行排名
func PublishHandler(c *gin.Context) {
	channel := c.Param("channel")
	infomation := c.PostForm("info")
	err := DB.Publish(c.Request.Context(), channel, infomation).Err()
	if err != nil {
		c.JSON(400, gin.H{"msg": err.Error()})
		return
	}
	c.JSON(200, gin.H{"info": "Publishing successfully"})
} //订阅发布者发布信息
func SubscribeHandler(c *gin.Context) {
	channel := c.Param("channel")
	pubsub := DB.Subscribe(c.Request.Context(), channel) //用于订阅频道的对象（channel）
	defer func(pubsub *redis.PubSub) {                   //在不用时或者出错时关闭（close）
		err := pubsub.Close()
		if err != nil {
			c.JSON(400, gin.H{"msg": err.Error()})
			return
		}
	}(pubsub)
	messageChannel := pubsub.Channel()
	for message := range messageChannel { //循环获取频道发布的消息
		fmt.Println("Received massage:", message.Payload) //输出在运行行中
		c.SSEvent("message", message.Payload)             //输出在运行行中//接收消息要用前端？(javascript)?
	}
}
