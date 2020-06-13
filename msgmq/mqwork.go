package msgmq

import (
	"github.com/golang/glog"
	"github.com/streadway/amqp"
	"log"
)

// rabbitmq 工作队列模式

// 发送消息(生产者)
func PublishMsgWork(queueName string,msg string) {
	//申请队列,如果不存在会自动创建，存在跳过创建，保证队列存在，消息能发送到队列中
	channel, err1 := MQInstance().GetMQChannel()
	if err1 != nil{
		glog.Fatal(err1)
	}
	_,err := channel.QueueDeclare(
		queueName,
		//控制队列是否为持久的，当mq重启的时候不会丢失队列
		true,
		//是否为自动删除
		false,
		//是否具有排他性
		false,
		//是否阻塞
		false,
		//额外属性
		nil,
	)
	if err != nil{
		glog.Fatal(err)
	}

	//发送消息到队列中
	_ = channel.Publish(
		"",
		queueName,
		//如果为true，根据exchange类型和routekey类型，如果无法找到符合条件的队列，name会把发送的信息返回给发送者
		false,
		//如果为true，当exchange发送到消息队列后发现队列上没有绑定的消费者,则会将消息返还给发送者
		false,
		//发送信息
		amqp.Publishing{
			// 设置消息为持久的
			DeliveryMode:    amqp.Persistent,
			ContentType:     "text/plain",
			Body:            []byte(msg),
		},
	)
	log.Printf("工作队列发送消息----- [x] %s", msg)
	defer channel.Close()
}

// 消费消息(消费者)
func ConsumeMsgWork(queueName string){
	//申请队列,如果不存在会自动创建，存在跳过创建，保证队列存在，消息能发送到队列中
	channel, _ := MQInstance().GetMQChannel()
	_,err := channel.QueueDeclare(
		queueName,
		//控制消息是否持久化，true开启
		true,
		//是否为自动删除
		false,
		//是否具有排他性
		false,
		//是否阻塞
		false,
		//额外属性
		nil,
	)
	if err != nil{
		glog.Fatal(err)
	}

	//接收消息
	msgs,err := channel.Consume(
		queueName,
		//用来区分多个消费者
		"",
		//是否自动应答(自动应答确认消息，这里设置为否，在下面手动应答确认)
		false,
		//是否具有排他性
		false,
		//如果设置为true，表示不能将同一个connection中发送的消息
		//传递给同一个connection的消费者
		false,
		//是否为阻塞
		false,
		nil,
	)
	if err != nil {
		glog.Fatal(err)
	}
	forever := make(chan bool)
	//启用协程处理消息
	go func() {
		for d := range msgs {
			// 手动确认收到本条消息, true表示回复当前信道所有未回复的ack，用于批量确认。false表示回复当前条目
			d.Ack(false)
			log.Printf("工作队列接收到消息----- [x] %s", d.Body)
			// todo 业务处理
		}
	}()
	<-forever
}