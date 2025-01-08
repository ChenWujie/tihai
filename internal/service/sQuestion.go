package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/goccy/go-json"
	"reflect"
	"strconv"
	"strings"
	"tihai/global"
	"tihai/internal/model"
	"tihai/utils"
	"time"
)

func CreateQuestion(question *model.Question) error {
	res := global.Db.Create(&question)
	if res.Error != nil {
		return res.Error
	}
	// TODO 解耦
	err := addQuestionToIndex(*question)
	if err != nil {
		return err
	}
	return nil
}

func UpdateQuestion(question model.Question) error {
	if err := global.Db.Model(&question).Omit("CreateAt").Updates(question).Error; err != nil {
		return err
	}
	err := updateQuestionByID(question)
	if err != nil {
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
	mapping := `{
	  "mappings": {
		"properties": {
		  "title": {
			"type": "text",
			"analyzer": "standard"
		  },
		  "content": {
			"type": "text",
			"analyzer": "standard"
		  },
		  "type": {
			"type": "keyword"
		  },
		  "image_url": {
			"type": "keyword",
			"index": false
		  },
		  "teacher_id": {
			"type": "integer",
			"index": false
		  },
		  "answer": {
			"type": "text",
			"index": false
		  }
		}
	  }
	}`

	req := esapi.IndicesCreateRequest{
		Index: "questions",
		Body:  strings.NewReader(mapping),
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
	questionJSON, err := json.Marshal(question)
	if err != nil {
		return fmt.Errorf("error marshaling article to JSON: %s", err)
	}

	req := esapi.IndexRequest{
		Index:      "questions",
		Body:       bytes.NewReader(questionJSON),
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

// SearchArticles 搜索文章
func SearchArticles(query string) (map[string]interface{}, error) {
	searchQuery := fmt.Sprintf(`{
		"query": {
		"bool": {
			"should": [
				{
					"multi_match": {
						"query": "%s",
						"fields": ["title", "content"]
					}
				},
				{
					"term": {
						"type": {
							"value": "%s"
						}
					}
				}
			]
		}
	}
	}`, query, query)

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

func updateQuestionByID(question model.Question) error {
	nonEmptyFields := make(map[string]interface{})
	// 通过反射获取结构体的类型和值
	structType := reflect.TypeOf(question)
	structValue := reflect.ValueOf(question)
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		value := structValue.Field(i).Interface()
		// 判断字段是否为空，根据不同类型判断空值情况
		if isNotEmpty(value) {
			nonEmptyFields[field.Tag.Get("json")] = value
		}
	}

	nonEmptyFieldsJSON, err := json.Marshal(map[string]interface{}{"doc": nonEmptyFields})
	if err != nil {
		return fmt.Errorf("error marshaling non-empty fields to JSON: %s", err)
	}

	updateRequest := esapi.UpdateRequest{
		Index:      "questions",
		DocumentID: strconv.Itoa(int(question.ID)),
		Body:       bytes.NewReader(nonEmptyFieldsJSON),
	}

	res, err := updateRequest.Do(context.Background(), global.ES)
	if err != nil {
		return fmt.Errorf("error updating question by ID: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error in updating question: %s", res.Status())
	}

	fmt.Println("Question updated successfully")
	return nil
}

// isNotEmpty判断值是否为空，根据不同类型做不同判断
func isNotEmpty(value interface{}) bool {
	switch v := value.(type) {
	case string:
		return v != ""
	case *string:
		return v != nil && *v != ""
	case int:
		return v != 0
	case *int:
		return v != nil && *v != 0
	// 可以根据实际结构体中的其他类型继续添加判断逻辑，比如切片、结构体指针等类型
	case time.Time:
		return !v.IsZero()
	default:
		return false
	}
}
