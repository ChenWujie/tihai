package model

import "gorm.io/gorm"

type Question struct {
	gorm.Model
	Title       string          `gorm:"type:text; not null" json:"title"`
	Public      bool            `gorm:"not null;type:bool" json:"public"`
	Content     string          `gorm:"type:text" json:"content"`
	ImageUrl    string          `gorm:"type:text" json:"url"`
	TeacherID   uint            `gorm:"not_null" json:"tid"`
	Type        string          `gorm:"type:enum('chose', 'multi_chose', 'judge', 'saq'); not null; default:'saq'" json:"type"`
	Answer      string          `gorm:"type:text" json:"answer"`
	Teacher     User            `gorm:"foreignKey:TeacherID" json:"teacher"`
	UserAnswers []StudentAnswer `gorm:"foreignKey:QuestionID"`
	Papers      []Paper         `gorm:"many2many:paper_questions;"`
}
