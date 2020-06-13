package msgmq

import (
	"fmt"
	"github.com/golang/glog"
	"github.com/streadway/amqp"
	_const "reallog/common/const"
	"sync"
)

type RabbitMQ struct {
	conn *amqp.Connection
	channel *amqp.Channel
}

var instance *RabbitMQ
var once sync.Once

func MQInstance() *RabbitMQ {
	once.Do(func() {
		instance = &RabbitMQ{}
	})
	return instance
}

func (r *RabbitMQ)InitConnect() bool {
	fmt.Println("=====================初始化RabbitMQ连接=====================")
	var err error
	r.conn, err = amqp.Dial(_const.RABBITMQ_BROKER_URL)
	if err != nil{
		glog.Infof("failed to connect rabbitmq: %v",err)
		return false
	}

	return true
}

func (r *RabbitMQ)CloseMq() {
	_ = r.channel.Close()
	_ = r.conn.Close()
}

func (r *RabbitMQ)GetMQCon() *amqp.Connection {
	return r.conn
}

func (r *RabbitMQ)GetMQChannel() (*amqp.Channel ,error){
	var err error
	r.channel, err = r.conn.Channel()
	if err != nil{
		glog.Infof("failed to open a rabbitmq channel: %v",err)
		return nil,err
	}
	return r.channel,nil
}
