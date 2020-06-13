package msgmq

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"reallog/common/vo"
	"strconv"
	"time"
	"unsafe"
)

// rabbitmq 路由routing模式

// 发送消息(生产者)
func PublishMsgRout(exName string,rtKey string, msg string) {

	ch, err := MQInstance().GetMQChannel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		exName, // name
		"direct",      // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	failOnError(err, "Failed to declare an exchange")
	err = ch.Publish(
		exName,          // exchange
		rtKey, // routing key
		//如果为true，根据exchange类型和routekey类型，如果无法找到符合条件的队列，name会把发送的信息返回给发送者
		true, // mandatory
		false, // immediate
		amqp.Publishing{
			// 消息持久化
			DeliveryMode: amqp.Persistent,
			ContentType: "text/plain",
			Body:        []byte(msg),
		})
	failOnError(err, "Failed to publish a message")
}


type Callback func(trainId string,msg string)

// 消费消息(消费者)
func ConsumeMsgRout(exName string,queName string,rtKey string,callback Callback) {
	ch, err := MQInstance().GetMQChannel()
	failOnError(err, "Failed to receive a message")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		exName,   // name
		"direct", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	failOnError(err, "Failed to receive a message")

	_, err = ch.QueueDeclare(
		queName,    // name
		false, // durable
		false, // delete when usused
		false,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to receive a message")

	err = ch.QueueBind(
		queName, // queue name
		rtKey,     // routing key
		exName, // exchange
		false,
		nil,
	)
	failOnError(err, "Failed to receive a message")

	msgs, err := ch.Consume(
		queName, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)

	forever := make(chan bool)
	go func() {
		for d := range msgs {
			var msgTrainLog vo.TrainLog
			_ = json.Unmarshal(d.Body, &msgTrainLog)
			trainId := msgTrainLog.TrainId
			id64 := strconv.FormatInt(trainId,10)

			// 避免循环引用  callback = ws.SendLogsToWeb(trainId string,msg string)
			strAddress := &callback
			strPointer := fmt.Sprintf("%d", unsafe.Pointer(strAddress))
			intPointer, _ := strconv.ParseInt(strPointer, 10, 0)
			var pointer *Callback
			pointer = *(**Callback)(unsafe.Pointer(&intPointer))
			(Callback)(*pointer)(id64,time.Now().String() + "====="+ msgTrainLog.TrainLog)

			log.Printf("路由队列接收到消息======== [x] %s", d.Body)
		}
	}()
	<-forever
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}