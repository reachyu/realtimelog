# 简介
* go+RabbitMQ+websocket实现在web端实时展示日志以及服务异步消息
# 目录介绍
## common
公共方法
## public
静态文件目录
* ws1.html和ws2.html是RabbitMQ 路由模式验证前端
* f1.html和f2.html是RabbitMQ 发布订阅模式验证前端
* k1.html和k2.html是kafka验证前端
### httpserver
http服务
### msgmq
消息队列
### ws
websocket代码目录
# 编译
```
go build
```

# 启动
```
go run main.go
```
