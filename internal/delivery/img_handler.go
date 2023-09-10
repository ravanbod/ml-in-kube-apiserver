package delivery

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type imgHandler struct {
	redisConn *redis.Client
}

type imageRequest struct {
	FileUrl string `json:"fileUrl" binding:"required"`
}

func (r *Handler) SetImgHandler(redisConn *redis.Client) {
	myImgHandler := imgHandler{redisConn: redisConn}

	r.Engine.GET("/", myImgHandler.index)
	r.Engine.POST("/image/", myImgHandler.addFileToQueue)
}

func (r *imgHandler) index(c *gin.Context) {
	c.JSON(200, "Hello, World!")
}

func (r *imgHandler) addFileToQueue(c *gin.Context) {
	var req imageRequest
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(400, "Bad request")
		return
	}
	r.redisConn.HSet(context.Background(), uuid.NewString(), "to_be_processed_url", req.FileUrl, "processed_url", "")
	c.JSON(200, "Added to queue")
}
