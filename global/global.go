package global

import (
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/go-redis/redis"
	"github.com/gorilla/websocket"
	"github.com/streadway/amqp"
	"gorm.io/gorm"
	"net/http"
	"sync"
)

var (
	Db      *gorm.DB
	RedisDB *redis.Client
	ES      *elasticsearch.Client
	Conn    *amqp.Connection
	Channel *amqp.Channel
)

var Upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	Conn  *websocket.Conn
	Mutex sync.Mutex
}

var UserClients map[uint]*Client
