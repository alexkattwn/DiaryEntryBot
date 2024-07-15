package bot

// Этот файл содержит основную логику работы бота

import (
	"DiaryEntryBot/config"
	"DiaryEntryBot/internal/repository"
	"DiaryEntryBot/internal/services"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// StartBot запускает бота
func StartBot(config *config.Config, db repository.Database) {
	// Создает нового бота, используя API-токен
	bot, err := tgbotapi.NewBotAPI(config.TelegramBotToken)
	if err != nil {
		log.Panic(err, "Не удалось создать бота")
	}

	log.Println("Бот запущен")

	//bot.Debug = true

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	// Создает канал для получения обновлений от Telegram
	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Panic(err, "Не удалось создать канал для получения обновлений от telegram")
	}

	service := services.NewDiaryService(db, config.EncryptionKey)

	// Обработка входящих сообщений
	for update := range updates {
		if update.Message != nil {
			handleMessage(update, bot, service)
		}
	}
}