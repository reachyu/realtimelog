package vo

// RabbitMQ 消息实体

type TrainLog struct {
	TrainId     int64   `json:"trainId"`
	TrainLog    string  `json:"trainLog"`
}
