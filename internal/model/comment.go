package model

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	AnswerId      int           `gorm:"not null"`
	TeacherId     int           `gorm:"not null" json:"tid"`
	StudentAnswer StudentAnswer `gorm:"foreignkey:AnswerID" json:"student_answer"`
	Teacher       User          `gorm:"foreignkey:TeacherID" json:"teacher"`
}
