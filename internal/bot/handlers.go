package bot

// Этот файл содержит обработчики сообщений бота

import (
	"DiaryEntryBot/internal/services"
	"fmt"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var (
	notificationTimer *time.Timer
)

// handleMessage обрабатывает входящие сообщения
func handleMessage(update tgbotapi.Update, bot *tgbotapi.BotAPI, service *services.DiaryService) {
	if update.Message.IsCommand() {
		switch update.Message.Command() {
		case "start":
			handleStartCommand(update, bot)
			return
		case "view":
			viewEntries(update, bot, service)
			return
		case "edit": 
			editEntry(update, bot, service)
			return
		case "delete":
			deleteEntry(update, bot, service)
			return
		}
	}

	userID := update.Message.From.ID
	content := update.Message.Text

	// Создание новой записи в дневнике
	service.CreateEntry(userID, content)

	resetTimer(update.Message.Chat.ID, bot)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Запись добавлена в дневник!")
	bot.Send(msg)
}

// viewEntries отображает записи дневника
func viewEntries(update tgbotapi.Update, bot *tgbotapi.BotAPI, service *services.DiaryService) {
	userID := update.Message.From.ID
	entries, err := service.GetEntries(userID)

	if err != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка получения записей.")
		bot.Send(msg)
		return
	}

	var response string
	for _, entry := range entries {
		response += fmt.Sprintf("ID: %d\nДата: %s\n%s\n\n", entry.ID, entry.CreatedAt.Format("02.01.2006 15:04"), entry.Content)
	}

	if response == "" {
		response = "У вас нет записей в дневнике."
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, response)
	bot.Send(msg)
}

// editEntry редактирует запись в дневнике
func editEntry(update tgbotapi.Update, bot *tgbotapi.BotAPI, service *services.DiaryService) {
	parts := strings.SplitN(update.Message.Text, " ", 3)
	if len(parts) < 3 {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Использование: /edit <ID> <новый текст>")
		bot.Send(msg)
		return
	}

	entryID, err := strconv.ParseUint(parts[1], 10, 32)
	if err != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Некорректный ID записи.")
		bot.Send(msg)
		return
	}

	content := parts[2]
	userID := update.Message.From.ID

	err = service.UpdateEntry(uint(entryID), userID, content)
	if err != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка при обновлении записи.")
		bot.Send(msg)
		return
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Запись обновлена!")
	bot.Send(msg)
}

// deleteEntry удаляет запись из дневника
func deleteEntry(update tgbotapi.Update, bot *tgbotapi.BotAPI, service *services.DiaryService) {
	parts := strings.SplitN(update.Message.Text, " ", 2)
	if len(parts) < 2 {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Использование: /delete <ID>")
		bot.Send(msg)
		return
	}

	entryID, err := strconv.ParseUint(parts[1], 10, 32)
	if err != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Некорректный ID записи.")
		bot.Send(msg)
		return
	}

	userID := update.Message.From.ID

	err = service.DeleteEntry(uint(entryID), userID)
	if err != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка при удалении записи.")
		bot.Send(msg)
		return
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Запись удалена!")
	bot.Send(msg)
}

// handleStartCommand обрабатывает команду /start
func handleStartCommand(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	welcomeMessage := "Привет! 😊 Я твой личный дневник-бот. Ты можешь отправить мне сообщение, и я сохраню его в твой дневник. " +
		"Ты также можешь использовать следующие команды:\n" +
		"/view - Просмотреть все записи\n" +
		"/edit <ID> <новый текст> - Редактировать запись\n" +
		"/delete <ID> - Удалить запись"

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, welcomeMessage)
	bot.Send(msg)
}

// sendNotification отправляет уведомление, если нет новых записей в течение 24 часов
func sendNotification(chatID int64, bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(chatID, "Давно не было новых записей! 📔")
	bot.Send(msg)

	resetTimer(chatID, bot)
}

func resetTimer(chatID int64, bot *tgbotapi.BotAPI) {
	if notificationTimer != nil {
		notificationTimer.Stop()
	}

	notificationTimer = time.AfterFunc(24 * time.Hour, func() {
		sendNotification(chatID, bot)
	})
}