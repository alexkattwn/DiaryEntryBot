package repository

// Этот файл содержит интерфейсы для работы с базой данных

import "DiaryEntryBot/internal/models"

type Database interface {
	CreateEntry(entry *models.DiaryEntry) error
	GetEntries(userID int) ([]models.DiaryEntry, error)
	UpdateEntry(entry *models.DiaryEntry) error
	DeleteEntry(entryID uint, userID int) error
	Close() error
}