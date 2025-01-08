package model

import (
	"gorm.io/gorm"
	"time"
)

type Paper struct {
	gorm.Model
	PaperName string     `gorm:"type:varchar(20);not null" json:"paper_name"`
	StartTime time.Time  `gorm:"column:start_time;type:datetime;not null" json:"start_time"`
	EndTime   time.Time  `gorm:"column:end_time;type:datetime;not null" json:"end_time"`
	UserID    uint       `gorm:"column:user_id;type:int(11);not null"`
	User      User       `gorm:"foreignKey:UserID" json:"user"`
	Questions []Question `gorm:"many2many:paper_questions;"`
	Classes   []Class    `gorm:"many2many:class_papers;"`
}
