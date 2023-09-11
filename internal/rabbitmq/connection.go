package rabbitmq

import (
	"ml-in-kube-apiserver/internal/config"

	"github.com/rabbitmq/amqp091-go"
)

func NewRabbitmeConnection(cfg config.RabbitmqConfig) (*amqp091.Connection, error) {
	return amqp091.Dial("amqp://" + cfg.User + ":" + cfg.Password + "@" + cfg.Host + "/")
}
