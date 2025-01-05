package global

import (
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

var (
	Db      *gorm.DB
	RedisDB *redis.Client
	ES      *elasticsearch.Client
)
