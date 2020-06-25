package _const

const (

    // 格式：amqp://用户名:密码@地址:端口号/host
    RABBITMQ_BROKER_URL = "amqp://test:1234@ip:5672/testvh"
	RABBITMQ_ROUT_EXCHANGE_NAME  = "ex.log01"
	RABBITMQ_FANOUT_EXCHANGE_NAME  = "ex.log02"
	RABBITMQ_ROUT_QUEUE  = "que.logque."
	RABBITMQ_ROUT_ROUTING_KEY    = "rtk.train."
	RABBITMQ_WORK_PUBLISH_QUEUE  = "mywork.que01"
	RABBITMQ_WORK_CONSUME_QUEUE  = "mywork.que01"
	RABBITMQ_FANOUT_QUEUE  = "que.fanout.logque."

	KAFKA_LOGS_TOPIC = "topic.logs."
	KAFKA_LOGS_GROUP = "group.logs."

	HTTP_SERVER_PORT = "9090"

)

var KAFKA_SERVER = []string{"ip1:9092","ip2:9092","ip3:9092"}
