package model

import "time"

type StudentAnswer struct {
	Id             int       `gorm:"primary_key;AUTO_INCREMENT" json:"id"`
	UserId         uint      `gorm:"not null" json:"user_id"`
	QuestionID     uint      `gorm:"not null" json:"question_id"`
	AnswerText     string    `json:"answer_text"`
	AnswerImageUrl string    `json:"url"`
	SubmitTime     time.Time `gorm:"not null" json:"submit_time"`
	User           User      `json:"student"`
}
