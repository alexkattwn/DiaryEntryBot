package main

// Этот файл содержит основной код для запуска бота

import (
	"log"

	"github.com/joho/godotenv"

	"DiaryEntryBot/config"
	"DiaryEntryBot/internal/bot"
	"DiaryEntryBot/internal/repository/postgres"
)

func main() {
	// загрузка переменных окружения
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Настройка конфигурации
	config := config.NewConfig()

	log.Println("Переменные окружения загружены")

	// Установка соединения с базой данных PostgreSQL
	db, err := postgres.NewPostgresDB(config)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Подключение к БД успешно")

	defer db.Close()

	// Запуск бота
	bot.StartBot(config, db)
}