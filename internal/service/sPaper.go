package service

import (
	"reflect"
	"tihai/global"
	"tihai/internal/model"
)

func CreatePaper(paper *model.Paper, questionIds []uint) error {
	tx := global.Db.Begin()

	if err := tx.Create(paper).Error; err != nil {
		tx.Rollback()
		return err
	}

	var questions []model.Question
	if err := tx.Where("id IN (?)", questionIds).Find(&questions).Error; err != nil {
		tx.Rollback()
		return err
	}

	for _, question := range questions {
		if err := tx.Model(paper).Association("Questions").Append(&question); err != nil {
			tx.Rollback()
			return err
		}
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func DeletePaper(paper model.Paper) error {
	tx := global.Db.Begin()

	if err := tx.Select("Questions").Delete(&paper).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func UpdatePaper(paper model.Paper, questionIds []uint) error {
	tx := global.Db.Begin()
	var questions []model.Question
	if err := tx.Where("id IN (?)", questionIds).Find(&questions).Error; err != nil {
		tx.Rollback()
		return err
	}
	err := tx.Model(&paper).Association("Questions").Replace(questions)
	if err != nil {
		tx.Rollback()
		return err
	}
	nonEmptyFields := make(map[string]interface{})
	// 通过反射获取结构体的类型和值
	structType := reflect.TypeOf(paper)
	structValue := reflect.ValueOf(paper)
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		value := structValue.Field(i).Interface()
		// 判断字段是否为空，根据不同类型判断空值情况
		if isNotEmpty(value) {
			nonEmptyFields[field.Tag.Get("json")] = value
		}
	}
	if err := tx.Model(&paper).Where("id = ?", paper.ID).Updates(nonEmptyFields).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func GetPaper(uid uint) ([]model.Paper, error) {
	var paper []model.Paper
	if err := global.Db.Preload("Questions").Find(&paper, "user_id = ?", uid).Error; err != nil {
		return paper, err
	}
	return paper, nil
}
