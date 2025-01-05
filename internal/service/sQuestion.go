package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/goccy/go-json"
	"io"
	"strconv"
	"strings"
	"tihai/global"
	"tihai/internal/model"
	"tihai/utils"
)

func CreateQuestion(question model.Question) error {
	if err := global.Db.AutoMigrate(&question); err != nil {
		return err
	}
	res := global.Db.Create(&question)
	if res.Error != nil {
		return res.Error
	}
	// TODO 解耦
	err := addQuestionToIndex(question)
	if err != nil {
		return err
	}
	return nil
}

func UpdateQuestion(question model.Question) error {
	if err := global.Db.Model(&question).Omit("CreateAt").Updates(question).Error; err != nil {
		return err
	}
	return nil
}

func DeleteQuestion(question model.Question) error {
	if err := global.Db.Delete(&question).Error; err != nil {
		return err
	}
	// TODO 解耦
	err := deleteQuestionByID(strconv.Itoa(int(question.ID)))
	if err != nil {
		return err
	}
	return nil
}

func FindList(t string, token string) ([]model.Question, error) {
	var question []model.Question
	_, err := global.RedisDB.Get(token).Result()
	if err == nil { // token失效
		return FindListByGuest(t)
	}
	//解析token，获取uid
	authMap, err := utils.ParseJWT(token)
	if err != nil {
		return nil, err
	}
	if err := global.Db.Where("(type = ? and public = ?) OR (type = ? and public = ? and teacher_id = ?) ", t, true, t, false, authMap["uid"]).Find(&question).Error; err != nil {
		return nil, err
	}
	return question, nil
}

func FindListByGuest(t string) ([]model.Question, error) {
	var question []model.Question
	if err := global.Db.Where("type = ? and public = ?", t, true).Find(&question).Error; err != nil {
		return nil, err
	}
	return question, nil
}

func LikeQuestion(uid, qid uint) (int64, string, error) {
	questionKey := fmt.Sprintf("question:like_count:%d", qid)
	userKey := fmt.Sprintf("user:%d:liked_articles", uid)
	global.RedisDB.SetNX(questionKey, 0, 0)
	// 先判断是否是已收藏，如果是，取消收藏
	isLiked, err := global.RedisDB.SIsMember(userKey, qid).Result()
	if err != nil {
		return 0, "", err
	}
	if isLiked {
		result, err := global.RedisDB.Decr(questionKey).Result()
		if err != nil {
			return 0, "", errors.New("取消收藏失败")
		}
		global.RedisDB.SRem(userKey, qid)
		return result, "取消收藏", nil
	}
	// 收藏文章
	result, err := global.RedisDB.Incr(questionKey).Result()
	if err != nil {
		return 0, "", err
	}
	global.RedisDB.SAdd(userKey, qid)
	return result, "收藏成功", nil
}

// CreateQuestionIndexWithMapping 创建问题索引
func CreateQuestionIndexWithMapping() error {
	mapping := map[string]interface{}{
		"mappings": map[string]interface{}{
			"properties": map[string]interface{}{
				"title": map[string]interface{}{
					"type":     "text",
					"analyzer": "standard",
				},
				"content": map[string]interface{}{
					"type":     "text",
					"analyzer": "standard",
				},
				"type": map[string]interface{}{
					"type": "keyword",
				},
				"image_url": map[string]interface{}{
					"type":  "keyword",
					"index": false,
				},
				"teacher_id": map[string]interface{}{
					"type":  "integer",
					"index": false,
				},
				"answer": map[string]interface{}{
					"type":  "text",
					"index": false,
				},
			},
		},
	}

	mappingJSON, err := json.Marshal(mapping)
	if err != nil {
		return fmt.Errorf("error marshaling mapping: %s", err)
	}

	req := esapi.IndicesCreateRequest{
		Index: "questions",
		Body:  strings.NewReader(string(mappingJSON)),
	}

	res, err := req.Do(context.Background(), global.ES)
	if err != nil {
		return fmt.Errorf("error creating index with mapping: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error creating index: %s", res.Status())
	}

	fmt.Println("Index created with custom mapping")
	return nil
}

// addArticleToIndex将文章添加到Elasticsearch的questions索引中
func addQuestionToIndex(question model.Question) error {
	// 将文章结构体转换为JSON格式的字节切片
	articleJSON, err := json.Marshal(question)
	if err != nil {
		return fmt.Errorf("error marshaling article to JSON: %s", err)
	}

	req := esapi.IndexRequest{
		Index:      "questions",
		Body:       io.NopCloser(strings.NewReader(string(articleJSON))),
		DocumentID: strconv.Itoa(int(question.ID)),
	}

	res, err := req.Do(context.Background(), global.ES)
	if err != nil {
		return fmt.Errorf("error adding article to index: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error adding article: %s", res.Status())
	}
	return nil
}

// 搜索文章
func SearchArticles(query string) (map[string]interface{}, error) {
	searchQuery := fmt.Sprintf(`{
		"query": {
			"multi_match": {
				"query": "%s",
				"fields": ["title", "content", "category"]
			}
		}
	}`, query)

	req := esapi.SearchRequest{
		Index: []string{"questions"},
		Body:  strings.NewReader(searchQuery),
	}

	res, err := req.Do(context.Background(), global.ES)
	if err != nil {
		return nil, fmt.Errorf("error searching documents: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("error searching documents: %s", res.Status())
	}

	fmt.Println("Search Results:")
	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, fmt.Errorf("error parsing the response body: %s", err)
	}
	return r, nil
}

func deleteQuestionByID(id string) error {
	deleteRequest := esapi.DeleteRequest{
		Index:      "questions",
		DocumentID: id,
	}

	res, err := deleteRequest.Do(context.Background(), global.ES)
	if err != nil {
		return fmt.Errorf("error deleting question by ID: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error in deleting question: %s", res.Status())
	}

	return nil
}
