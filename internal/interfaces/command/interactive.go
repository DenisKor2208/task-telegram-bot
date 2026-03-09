package command

import (
	"context"

	"github.com/DenisKor2208/task-telegram-bot/internal/interfaces/messaging"
)

type InteractiveCommand interface {
	Command
	//NextStep - Следующая команда после ввода
	NextStep() Command // следующая команда после ввода

	//InputField - Ключ для сохранения ввода в сессию
	InputField() string // куда сохранить ввод

	//Prompt
	Prompt(context.Context, Model, messaging.Message) error
}
