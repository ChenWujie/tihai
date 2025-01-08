package model

import "gorm.io/gorm"

type Class struct {
	gorm.Model
	AdminID   uint    `json:"admin_id"`
	ClassName string  `json:"class_name"`
	Admin     User    `gorm:"foreignKey:AdminID" json:"admin"`
	Users     []User  `gorm:"many2many:class_users;" json:"users"`
	Papers    []Paper `gorm:"many2many:class_papers" json:"papers"`
}
