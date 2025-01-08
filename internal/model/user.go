package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string  `gorm:"type:varchar(20);unique;not null" json:"username"`
	Password string  `gorm:"type:varchar(70);not null" json:"password"`
	Role     string  `gorm:"type:enum('student', 'teacher', 'admin');not null;default:'student'" json:"role"`
	Email    string  `gorm:"type:varchar(30)" json:"email"`
	Nickname string  `gorm:"size:16" json:"nickname"`
	Classes  []Class `gorm:"many2many:class_users;"`
}
