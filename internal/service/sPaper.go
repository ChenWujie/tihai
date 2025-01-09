package service

import (
	"errors"
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

func AssignPapers(uid, paperId uint, classIds []uint) error {
	tx := global.Db.Begin()
	var paper model.Paper
	if err := tx.First(&paper, paperId).Error; err != nil {
		tx.Rollback()
		return err
	}
	if paper.UserID != uid {
		tx.Rollback()
		return errors.New("仅试卷创建者可分配试卷！")
	}
	classes := make([]model.Class, 0)
	for _, classId := range classIds {
		var class model.Class
		if err := tx.First(&class, classId).Error; err != nil {
			tx.Rollback()
			return err
		}
		classes = append(classes, class)
	}
	err := tx.Model(&paper).Association("Classes").Append(&classes)
	if err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func QueryClassPapers(uid uint) ([]model.Paper, error) {
	var user model.User
	user.ID = uid
	var classes []model.Class
	err := global.Db.Model(&user).Association("Classes").Find(&classes)
	if err != nil {
		return nil, err
	}
	var papers []model.Paper
	for _, v := range classes {
		classPapers := make([]model.Paper, 0)
		err := global.Db.Model(&v).Association("Papers").Find(&classPapers)
		if err != nil {
			return nil, err
		}
		papers = append(papers, classPapers...)
	}
	return papers, nil
}
