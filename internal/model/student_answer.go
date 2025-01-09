package model

import (
	"gorm.io/gorm"
	"time"
)

type StudentAnswer struct {
	gorm.Model
	UserID         uint      `gorm:"not null" json:"user_id"`
	QuestionID     uint      `gorm:"not null" json:"question_id"`
	AnswerText     string    `json:"answer_text"`
	AnswerImageUrl string    `json:"url"`
	SubmitTime     time.Time `gorm:"not null" json:"submit_time"`
	User           User      `gorm:"foreignkey:UserID" json:"student"`
}
