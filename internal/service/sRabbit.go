package service

import (
	"fmt"
	"github.com/goccy/go-json"
	"github.com/streadway/amqp"
	"log"
	"tihai/global"
	"tihai/internal/model"
)

func PublishMessage(data []byte) error {
	err := global.Channel.Publish("paper", "note", false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        data,
	})
	if err != nil {
		log.Printf("failed to publish a message: %v", err)
	} else {
		log.Printf("message published successfully")
	}
	return err
}

func ConsumeMessage() {
	msgs, err := global.Channel.Consume(
		"notification", // 队列名称
		"",             // 消费者名称
		true,           // 自动确认
		false,          // 独占
		false,          // 无本地
		false,          // 无等待
		nil,            // 参数
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			var paper model.Paper
			err := json.Unmarshal(d.Body, &paper)
			if err != nil {
				log.Printf("failed to unmarshal the paper: %v", err)
			}
			message := fmt.Sprintf("新试卷发布啦！%s ，请在%s~%s完成", paper.PaperName, paper.StartTime, paper.EndTime)
			// 推送给班级
			var classes []model.Class
			err = global.Db.Model(&paper).Association("Classes").Find(&classes)
			if err != nil {
				log.Printf("Failed to get classes: %v", err)
				return
			}
			for _, class := range classes {
				var users []model.User
				err := global.Db.Model(&class).Association("Users").Find(&users)
				if err != nil {
					log.Printf("Failed to get users: %v", err)
				}

				for _, user := range users {
					// 推送给每个用户
					if client, ok := global.UserClients[user.ID]; ok {
						client.Mutex.Lock()
						err := client.Conn.WriteMessage(1, []byte(message))
						client.Mutex.Unlock()
						if err != nil {
							log.Println(err)
						}
					}
				}
				if client, ok := global.UserClients[class.AdminID]; ok {
					client.Mutex.Lock()
					err := client.Conn.WriteMessage(1, []byte(message))
					client.Mutex.Unlock()
					if err != nil {
						log.Println(err)
					}
				}
			}
		}
	}()

	log.Printf("Waiting for messages. To exit press CTRL+C")
	<-forever
}
