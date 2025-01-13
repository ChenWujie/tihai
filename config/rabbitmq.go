package config

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"tihai/global"
)

func initRabbitMQ() {
	addr := AppConfig.Rabbit.Host
	port := AppConfig.Rabbit.Port
	username := AppConfig.Rabbit.Username
	password := AppConfig.Rabbit.Password
	virtualHost := AppConfig.Rabbit.VirtualHost
	if virtualHost == "" {
		virtualHost = "/" // 设置默认虚拟主机值
	}
	// 构建连接字符串
	connStr := fmt.Sprintf("amqp://%s:%s@%s:%d/%s", username, password, addr, port, virtualHost)

	// 连接到RabbitMQ服务器
	conn, err := amqp.Dial(connStr)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ, got error: %v", err)
	}

	// 创建通道
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to create channel, got error: %v", err)
	}

	// 声明队列
	queue, err := ch.QueueDeclare(
		"notification",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare queue, got error: %v", err)
	}

	// 声明交换机
	err = ch.ExchangeDeclare(
		"paper",
		"direct",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare exchange, got error: %v",
			err)
	}

	// 绑定队列和交换机
	err = ch.QueueBind(
		queue.Name,
		"note",
		"paper",
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to bind queue and exchange, got error: %v",
			err)
	}

	// 将连接实例赋值给全局变量，方便其他地方使用
	global.Conn = conn
	global.Channel = ch
}
