package service

import (
	"errors"
	"strconv"
	"tihai/global"
	"tihai/internal/model"
)

func CreateStudentAnswer(answer model.StudentAnswer) (right bool, rate float64, err error) {
	// 回答
	if err := global.Db.Create(&answer).Error; err != nil {
		return false, 0, err
	}
	subkey := strconv.Itoa(int(answer.UserID)) + ":" + strconv.Itoa(int(answer.QuestionID))
	if count, err := global.RedisDB.Exists(subkey).Result(); err != nil {
		return false, 0, err
	} else {
		if count == 0 {
			global.RedisDB.Set(subkey, "1", 0)
		}
	}
	// 客观题，返回答案
	var question model.Question
	if err := global.Db.Model(&question).First(&question, answer.QuestionID).Error; err != nil {
		return false, 0, err
	}
	if question.Type != "saq" {
		keyright := strconv.Itoa(int(answer.QuestionID)) + ":right"
		keyfalse := strconv.Itoa(int(answer.QuestionID)) + ":false"
		global.RedisDB.SetNX(keyright, 0, 0)
		global.RedisDB.SetNX(keyfalse, 0, 0)
		var flag bool
		if answer.AnswerText == question.Answer {
			flag = true
			global.RedisDB.Incr(keyright)
		} else {
			flag = false
			global.RedisDB.Incr(keyfalse)
		}
		res, _ := global.RedisDB.Get(keyright).Result()
		right, _ := strconv.Atoi(res)
		res, _ = global.RedisDB.Get(keyfalse).Result()
		fault, _ := strconv.Atoi(res)
		return flag, float64(right) / float64(right+fault), err
	}
	return false, -1, nil
}

func GetStudentAnswerList(uid uint) ([]model.StudentAnswer, error) {
	var answers []model.StudentAnswer
	err := global.Db.Where("user_id = ?", uid).Find(&answers).Error
	return answers, err
}

func GetUserAnswerList(userId, questionId uint) ([]model.StudentAnswer, error) {
	subkey := strconv.Itoa(int(userId)) + ":" + strconv.Itoa(int(questionId))
	if count, err := global.RedisDB.Exists(subkey).Result(); err != nil {
		return nil, err
	} else {
		if count == 0 {
			return nil, errors.New("提交答案后可查看")
		}
	}
	var answers []model.StudentAnswer
	if err := global.Db.Where("question_id = ?", questionId).Find(&answers).Error; err != nil {
		return nil, err
	}
	return answers, nil
}
