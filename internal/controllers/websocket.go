package controllers

import (
	"github.com/gin-gonic/gin"
	"log"
	"tihai/global"
)

func WsHandler(c *gin.Context) {
	ws, err := global.Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	uid, _ := c.Get("uid")
	client := &global.Client{Conn: ws}
	global.UserClients[uid.(uint)] = client
	log.Println("创建WS连接")
}
