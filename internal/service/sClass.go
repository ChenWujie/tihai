package service

import (
	"errors"
	"gorm.io/gorm"
	"tihai/global"
	"tihai/internal/model"
)

func CreateClass(class model.Class) (uint, error) {
	tx := global.Db.Begin()
	if err := tx.Create(&class).Error; err != nil {
		tx.Rollback()
		return 0, err
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return 0, err
	}
	return class.ID, nil
}

func JoinClass(class model.Class, uid uint) error {
	tx := global.Db.Begin()
	if err := tx.First(&class, class.ID).Error; err != nil {
		tx.Rollback()
		return err
	}
	if class.AdminID == uid {
		tx.Rollback()
		return errors.New("创建者不用加入班级！")
	}
	var user model.User
	if err := tx.First(&user, uid).Error; err != nil {
		tx.Rollback()
		return err
	}
	err := tx.Model(&class).Association("Users").Append(&user)
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

func DeleteClass(class model.Class, uid uint) error {
	tx := global.Db.Begin()
	if err := tx.First(&class, class.ID).Error; err != nil {
		tx.Rollback()
		return err
	}
	if class.AdminID != uid {
		tx.Rollback()
		return errors.New("非法操作")
	}
	if err := tx.Select("Users").Delete(&class).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func QueryClass(user model.User) ([]model.Class, error) {
	classes := make([]model.Class, 0)
	err := global.Db.Model(&user).Association("Classes").Find(&classes)
	if err != nil {
		return nil, err
	}
	var createdClass []model.Class
	if err := global.Db.Model(&model.Class{}).Preload("Users", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, username")
	}).Preload("Admin", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, username")
	}).Where("admin_id", user.ID).Find(&createdClass).Error; err != nil {
		return nil, err
	}
	for _, v := range createdClass {
		classes = append(classes, v)
	}
	return classes, nil
}
