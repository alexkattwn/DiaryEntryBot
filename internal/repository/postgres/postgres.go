package postgres

// Этот файл реализует интерфейс Database для PostgreSQL

import (
	"DiaryEntryBot/config"
	"DiaryEntryBot/internal/models"
	"DiaryEntryBot/internal/repository"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// PostgresDB представляет соединение с базой данных PostgreSQL
type PostgresDB struct {
	*gorm.DB
}

// NewPostgresDB создает новое соединение с базой данных PostgreSQL
func NewPostgresDB(config *config.Config) (repository.Database, error) {
	db, err := gorm.Open("postgres", config.PostgresDSN)

	if err != nil {
		return nil, err
	}

	// Автоматическое создание таблиц, если их нет
	err = db.AutoMigrate(&models.DiaryEntry{}).Error

	if err != nil {
		return nil, err
	}

	return &PostgresDB{db}, nil
}

// CreateEntry создает новую запись в дневнике
func (db *PostgresDB) CreateEntry(entry *models.DiaryEntry) error {
	return db.Create(entry).Error
}

// GetEntries получает все записи дневника для указанного пользователя
func (db *PostgresDB) GetEntries(userID int) ([]models.DiaryEntry, error) {
	var entries []models.DiaryEntry

	err := db.Where("user_id = ?", userID).Find(&entries).Error

	return entries, err
}

// UpdateEntry обновляет существующую запись в дневнике
func (db *PostgresDB) UpdateEntry(entry *models.DiaryEntry) error {
	return db.Save(entry).Error
}

// DeleteEntry удаляет запись из дневника
func (db *PostgresDB) DeleteEntry(entryID uint, userID int) error {
	return db.Where("id = ? AND user_id = ?", entryID, userID).Delete(&models.DiaryEntry{}).Error
}


// Close закрывает соединение с базой данных
func (db *PostgresDB) Close() error {
	return db.DB.Close()
}