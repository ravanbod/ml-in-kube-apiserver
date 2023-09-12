package delivery

import (
	"context"
	"encoding/json"
	"ml-in-kube-apiserver/internal/config"
	"net/http"
	"net/url"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	minio "github.com/minio/minio-go/v7"
	credentials "github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
)

type imgHandler struct {
	redisConn  *redis.Client
	rabbitConn *amqp091.Connection
	s3Config   config.S3Config
}

type imageRequest struct {
	FileUrl string `json:"url" binding:"required"`
}

type imageUpdateRequest struct {
	FileUrl string `json:"url" binding:"required"`
	FileId  string `json:"id" binding:"required"`
}

func (r *Handler) SetImgHandler(redisConn *redis.Client, rabbitConn *amqp091.Connection, s3Cfg config.S3Config) {
	myImgHandler := imgHandler{redisConn: redisConn, rabbitConn: rabbitConn, s3Config: s3Cfg}

	r.Engine.GET("/", myImgHandler.index)
	r.Engine.POST("/image/", myImgHandler.addFileToQueue)
	r.Engine.PATCH("/image/", myImgHandler.updateFilePredictedUrl)
	r.Engine.POST("/upload/", myImgHandler.uploadFileToAddInQueue)
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
	c.JSON(201, "Added to queue")
}

func (r *imgHandler) updateFilePredictedUrl(c *gin.Context) {
	var req imageUpdateRequest
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(400, "Bad request")
		return
	}
	r.redisConn.HSet(context.Background(), req.FileId, "processed_url", req.FileUrl)
	c.JSON(200, "updated")
}

func (r *imgHandler) uploadFileToAddInQueue(c *gin.Context) {

	// Upload and save file
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, "Bad request")
		return
	}
	extension := filepath.Ext(file.Filename)
	newFileName := uuid.New().String() + extension
	savedFileAddress := "/tmp/upload/" + newFileName
	if err := c.SaveUploadedFile(file, savedFileAddress); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Unable to save the file",
		})
		return
	}

	// upload file to minio

	minioClient, err := minio.New(r.s3Config.Host, &minio.Options{
		Creds:  credentials.NewStaticV4(r.s3Config.AccessKey, r.s3Config.SecretKey, ""),
		Secure: false,
	})

	_, err = minioClient.FPutObject(context.Background(), "bucket1", newFileName, savedFileAddress, minio.PutObjectOptions{})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Unable to save the file in minio",
		})
	}

	url, err := minioClient.PresignedGetObject(context.Background(), "bucket1", newFileName, time.Hour*24, url.Values{})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Unable to get url from minio",
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Your file has been successfully uploaded.",
		"url":     url.String(),
	})
}
