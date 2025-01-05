package config

import (
	"github.com/elastic/go-elasticsearch/v8"
	"log"
	"tihai/global"
	"tihai/internal/service"
)

func initElasticSearch() {
	addr := AppConfig.ElasticSearch.Addr
	addresses := make([]string, 0)
	addresses = append(addresses, "http://"+addr)
	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{
			"http://localhost:9200", // Elasticsearch 地址
		},
	})
	if err != nil {
		log.Fatalf("Failed to connect to ElasticSearch, got error: %v", err)
	}
	global.ES = es
	exists, _ := global.ES.Indices.Exists([]string{"questions"})
	if exists.StatusCode != 200 {
		err := service.CreateQuestionIndexWithMapping()
		if err != nil {
			log.Fatal("Failed to create question index")
		}
	}
}
