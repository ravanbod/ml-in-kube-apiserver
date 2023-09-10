package delivery

import (
	"ml-in-kube-apiserver/internal/config"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Engine *gin.Engine
	cfg    config.HTTPServerConfig
}

func NewHandler(cfg config.HTTPServerConfig) Handler {
	engine := gin.Default()
	return Handler{cfg: cfg, Engine: engine}
}

func (r *Handler) StartServer() error {
	err := r.Engine.Run(r.cfg.Host + ":" + r.cfg.Port)
	return err
}
