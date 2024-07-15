package services

// Этот файл содержит бизнес-логику для работы с дневником

import (
	"DiaryEntryBot/internal/models"
	"DiaryEntryBot/internal/repository"
	"DiaryEntryBot/internal/utils"
)

// DiaryService представляет сервис для работы с дневником
type DiaryService struct {
	db  repository.Database
	key string
}

// NewDiaryService создает новый сервис для работы с дневником
func NewDiaryService(db repository.Database, key string) *DiaryService {
	return &DiaryService{db, key}
}

// CreateEntry создает новую запись в дневнике
func (s *DiaryService) CreateEntry(userID int, content string) error {
	encryptedContent, err := utils.Encrypt(s.key, content)
	if err != nil {
		return err
	}

	entry := &models.DiaryEntry{
		UserID:  userID,
		Content: encryptedContent,
	}

	return s.db.CreateEntry(entry)
}

// GetEntries получает все записи дневника для указанного пользователя
func (s *DiaryService) GetEntries(userID int) ([]models.DiaryEntry, error) {
	entries, err := s.db.GetEntries(userID)
	if err != nil {
		return nil, err
	}

	for i, entry := range entries {
		decryptedContent, err := utils.Decrypt(s.key, entry.Content)
		if err != nil {
			return nil, err
		}
		entries[i].Content = decryptedContent
	}

	return entries, nil
}

// UpdateEntry обновляет запись в дневнике
func (s *DiaryService) UpdateEntry(entryID uint, userID int, content string) error {
	encryptedContent, err := utils.Encrypt(s.key, content)
	if err != nil {
		return err
	}

	entry := &models.DiaryEntry{
		ID:      entryID,
		UserID:  userID,
		Content: encryptedContent,
	}

	return s.db.UpdateEntry(entry)
}

// DeleteEntry удаляет запись из дневника
func (s *DiaryService) DeleteEntry(entryID uint, userID int) error {
	return s.db.DeleteEntry(entryID, userID)
}