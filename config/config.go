package config

// Этот файл содержит конфигурационные параметры

import (
	"os"
)

type Config struct {
	TelegramBotToken string
	PostgresDSN      string
	EncryptionKey    string
}

// NewConfig создает новую конфигурацию, используя переменные окружения
func NewConfig() *Config {
	return &Config{
		TelegramBotToken: os.Getenv("TELEGRAM_BOT_TOKEN"),
		PostgresDSN:      os.Getenv("POSTGRES_DSN"),
		EncryptionKey:    os.Getenv("ENCRYPTION_KEY"),
	}
}