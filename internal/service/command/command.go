package command

import "gopkg.in/telebot.v4"

type Command interface {
	Register(b *telebot.Bot)
	Name() string
}
