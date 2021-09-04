package bot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

// NewInlineKeyboardButtonWithQueryCurrentChat create a tgbotapi.InlineKeyboardButton with SwitchInlineQueryCurrentChat
func NewInlineKeyboardButtonWithQueryCurrentChat(text, query string) tgbotapi.InlineKeyboardButton {
	return tgbotapi.InlineKeyboardButton{
		Text:                         text,
		SwitchInlineQueryCurrentChat: &query,
	}
}
