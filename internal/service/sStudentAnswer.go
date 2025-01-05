package service

import (
	"tihai/global"
	"tihai/internal/model"
)

func CreateStudentAnswer(answer model.StudentAnswer) (right bool, err error) {
	if err := global.Db.AutoMigrate(&answer); err != nil {
		return false, err
	}
	// 回答
	if err := global.Db.Create(&answer).Error; err != nil {
		return false, err
	}
	// 客观题，返回答案
	var question model.Question
	if err := global.Db.Model(&question).First(&question, answer.QuestionID).Error; err != nil {
		return false, err
	}
	if question.Type != "saq" {
		if answer.AnswerText == question.Answer {
			return true, nil
		}
	}
	return false, nil
}

func GetStudentAnswerList(uid uint) ([]model.StudentAnswer, error) {
	var answers []model.StudentAnswer
	err := global.Db.Where("user_id = ?", uid).Find(&answers).Error
	return answers, err
}

func GetUserAnswerList(questionId uint) ([]model.StudentAnswer, error) {
	var answers []model.StudentAnswer
	if err := global.Db.Where("question_id = ?", questionId).Find(&answers).Error; err != nil {
		return nil, err
	}
	return answers, nil
}
