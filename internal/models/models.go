package models

// Этот файл определяет структуру данных

import (
	"time"
)

type DiaryEntry struct {
	ID        uint `gorm:"primary_key"`
	UserID    int
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}