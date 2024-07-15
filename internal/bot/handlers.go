package bot

// Этот файл содержит обработчики сообщений бота

import (
	"DiaryEntryBot/internal/services"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// handleMessage обрабатывает входящие сообщения
func handleMessage(update tgbotapi.Update, bot *tgbotapi.BotAPI, service *services.DiaryService) {
	userID := update.Message.From.ID
	content := update.Message.Text

	// Команды для работы с дневником
	if strings.HasPrefix(content, "/view") {
		viewEntries(update, bot, service)
		return
	} else if strings.HasPrefix(content, "/edit") {
		editEntry(update, bot, service)
		return
	} else if strings.HasPrefix(content, "/delete") {
		deleteEntry(update, bot, service)
		return
	}

	// Создание новой записи в дневнике
	service.CreateEntry(userID, content)

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
		response += "ID: " + strconv.Itoa(int(entry.ID)) + "\n"
		response += entry.CreatedAt.Format("2006-01-02 15:04:05") + "\n"
		response += entry.Content + "\n\n"
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