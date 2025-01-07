package service

import (
	"tihai/global"
	"tihai/internal/model"
)

func CreatePaper(paper model.Paper, questionIds []uint) error {
	tx := global.Db.Begin()

	if err := tx.Create(&paper).Error; err != nil {
		tx.Rollback()
		return err
	}

	var questions []model.Question
	if err := tx.Where("id IN (?)", questionIds).Find(&questions).Error; err != nil {
		tx.Rollback()
		return err
	}

	for _, question := range questions {
		if err := tx.Model(&paper).Association("Questions").Append(&question); err != nil {
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
