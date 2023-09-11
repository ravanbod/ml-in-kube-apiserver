package delivery

import (
	"context"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
)

type imgHandler struct {
	redisConn  *redis.Client
	rabbitConn *amqp091.Connection
}

type imageRequest struct {
	FileUrl string `json:"fileUrl" binding:"required"`
}

func (r *Handler) SetImgHandler(redisConn *redis.Client, rabbitConn *amqp091.Connection) {
	myImgHandler := imgHandler{redisConn: redisConn, rabbitConn: rabbitConn}

	r.Engine.GET("/", myImgHandler.index)
	r.Engine.POST("/image/", myImgHandler.addFileToQueue)
}

func (r *imgHandler) index(c *gin.Context) {
	c.JSON(200, "Hello, World!")
}

// I know that this function and the whole of code is a real SHIT, because it has no standard architecture, no clean code and no good things
// I don't have enough time to refactor this now. So, Blame youself for reading this shit
func (r *imgHandler) addFileToQueue(c *gin.Context) {
	var req imageRequest
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(400, "Bad request")
		return
	}

	fileUUID := uuid.NewString()

	r.redisConn.HSet(context.Background(), fileUUID, "to_be_processed_url", req.FileUrl, "processed_url", "")
	ch, err := r.rabbitConn.Channel()
	if err != nil {
		c.JSON(500, "rabbit error")
		return
	}

	type ImageDataForRabbit struct {
		FileUrl string `json:"url"`
		FileId  string `json:"id"`
	}

	imgPayload := ImageDataForRabbit{FileUrl: req.FileUrl, FileId: fileUUID}
	jsonData, err := json.Marshal(imgPayload)
	if err != nil {
		c.JSON(500, "json error")
		return
	}

	ch.PublishWithContext(context.Background(), "", "q1", false, false, amqp091.Publishing{
		ContentType: "application/json",
		Body:        jsonData,
	})
	c.JSON(200, "Added to queue")
}
