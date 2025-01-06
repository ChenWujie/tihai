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
			log.Fatalf("Failed to create question index, got error: %v", err)
		}
	}

	// 以下内容，当索引存在时删除当前es索引，然后创建，不存在就直接创建
	// 这是为了方便找bug改的，原来是不存在索引就创建
	// 现在不想改，下面一段都注释了
	//if exists.StatusCode == 200 {
	//	req := esapi.IndicesDeleteRequest{Index: []string{"questions"}}
	//	res, err := req.Do(context.Background(), global.ES)
	//	if err != nil {
	//		fmt.Println("error creating index with mapping: %s", err)
	//	}
	//	defer res.Body.Close()
	//}
	//err = service.CreateQuestionIndexWithMapping()
	//if err != nil {
	//	log.Fatal("Failed to create question index")
	//}
}
