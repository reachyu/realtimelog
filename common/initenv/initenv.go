package initenv

import (
	"log"
	"os"
	mq "reallog/msgmq"
)

func InitEnv() {
	initRabbitMQ()
}

func initRabbitMQ() {
	issucc := mq.MQInstance().InitConnect()
	if !issucc {
		log.Println("init rabbitmq failure...")
		os.Exit(1)
	}
}