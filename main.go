package main

import (
	"encoding/json"
	"flag"
	_const "reallog/common/const"
	"reallog/common/initenv"
	"reallog/httpserver"
	mq "reallog/msgmq"
	"time"
)

func main() {
	// 不调用的话，glog会报错----ERROR: logging before flag.Parse:
	flag.Parse()
	initenv.InitEnv()
	go httpserver.InitHttpServer()

	// 测试实时日志,rabbitmq路由模式
	//testlogs()

	// 测试工作队列模式
	testwork()

}

func testlogs()  {
	msg := map[string]interface{}{
		"trainId":    123456,
		"trainLog":   "测试测试测试测试",
	}
	msgJson, _ := json.Marshal(msg)
	forever := make(chan bool)
	go func() {
		for {
			// 测试，所以id写死
			mq.PublishMsgRout(_const.RABBITMQ_ROUT_EXCHANGE_NAME,_const.RABBITMQ_ROUT_ROUTING_KEY+"123456",string(msgJson))
			time.Sleep(1 * time.Second)
		}
	}()
	<-forever
}

func testwork()  {

	msg := map[string]interface{}{
		"msgId": "",
		"objName": "trial",
		"opt": "start",
		"data": map[string]interface{}{
			"trialId": 123,
			"modelId": 666,
			"dataPath": "image/test/001509.jpg",
			"trialPara": map[string]interface{}{"score_threshold": "0.4"},
			"ossBucket": "file",
			"ossPath": "result/detection/test",
			"model": map[string]interface{}{
				"modelPath":  "model/yolov3.YYMnist.h5",
				"configPath": "result/detection/train-06-12/config.zip",
			},
		},
	}

	msgJson, _ := json.Marshal(msg)

	go func() {
		for {
			mq.PublishMsgWork(_const.RABBITMQ_WORK_PUBLISH_QUEUE,time.Now().String() +"======="+string(msgJson))
			time.Sleep(1 * time.Second)
		}
	}()

	mq.ConsumeMsgWork(_const.RABBITMQ_WORK_CONSUME_QUEUE)

}
