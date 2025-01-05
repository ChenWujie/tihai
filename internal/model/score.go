package model

import "time"

type Score struct {
	Id         int       `gorm:"primary_key;AUTO_INCREMENT"`
	StudentId  int       `gorm:"not null" json:"student_id"`
	QuestionId int       `gorm:"not null" json:"question_id"`
	Sc         int       `json:"score"`
	GradedBy   int       `gorm:"foreignKey:GradedBy; references:users.id" json:"graded_by"`
	GradedTime time.Time `json:"graded_time"`
	Student    User      `gorm:"foreignKey:StudentId" json:"student"`
	Question   Question  `gorm:"foreignKey:QuestionId" json:"question"`
}
