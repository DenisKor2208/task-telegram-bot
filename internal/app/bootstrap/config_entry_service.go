package bootstrap

import (
	"github.com/DenisKor2208/task-telegram-bot/internal/bot/config"
	"github.com/DenisKor2208/task-telegram-bot/internal/interfaces/service"
	"github.com/DenisKor2208/task-telegram-bot/internal/model/messages"
)

func ConfigEntryService(msgModel *messages.Model) (*config.ConfigEntryService, error) {
	initalConfig := map[string]service.ConfigEntry{
		"start":              {TargetCommand: "", CommandFields: []string{}},
		"view_tasks":         {TargetCommand: "", CommandFields: []string{}},
		"filter_all":         {TargetCommand: "", CommandFields: []string{}},
		"filter_closed":      {TargetCommand: "", CommandFields: []string{}},
		"filter_completed":   {TargetCommand: "", CommandFields: []string{}},
		"filter_in_progress": {TargetCommand: "", CommandFields: []string{}},
		"filter_overdue":     {TargetCommand: "", CommandFields: []string{}},
		"add_task": {
			TargetCommand: "save_task",
			CommandFields: []string{"title"},
		},
		"edit_task": {
			TargetCommand: "edit_task_by_id",
			CommandFields: []string{"title"},
		},
		"delete_task": {
			TargetCommand: "",
			CommandFields: []string{"title"},
		},
		"completed_task": {
			TargetCommand: "",
			CommandFields: []string{"title"},
		},
		"closed_task": {
			TargetCommand: "",
			CommandFields: []string{"title"},
		},
		"closed_task_by_id": {
			TargetCommand: "",
			CommandFields: []string{"task_id"},
		},
		"completed_task_by_id": {
			TargetCommand: "",
			CommandFields: []string{"task_id"},
		},
		"edit_task_by_id": {
			TargetCommand: "update_task",
			CommandFields: []string{"task_id", "title"},
		},
		"delete_task_by_id": {
			TargetCommand: "",
			CommandFields: []string{"task_id"},
		},
	}

	service, err := config.NewConfigEntryService(msgModel, initalConfig)
	if err != nil {
		return nil, err
	}

	msgModel.SetConfigEntryService(service)
	return service, nil
}
