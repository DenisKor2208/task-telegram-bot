package edit_task

import (
	"strconv"
	"time"

	"github.com/DenisKor2208/task-telegram-bot/internal/model/messages"
	"github.com/DenisKor2208/task-telegram-bot/internal/models"
	"github.com/DenisKor2208/task-telegram-bot/internal/resources"
	"github.com/DenisKor2208/task-telegram-bot/internal/types"
	btnupdatetask "github.com/DenisKor2208/task-telegram-bot/internal/ui/buttons/commands/edit_task"
	"github.com/pkg/errors"
)

// UpdateTaskCommand - структура команды /update_task
type UpdateTaskCommand struct{}

func (c *UpdateTaskCommand) Execute(s types.Model, msg types.Message, session *models.UserSession) error {

	var task *models.Task

	if session.Data == nil {
		return errors.New("no data in session: arguments not found")
	}

	// Получаем ID задачи из сессии
	// taskId, ok := session.Data["task_id"].(string)
	taskId, ok := session.Data["payload_task_id"].(int64)
	if !ok {
		// Fallback на старый способ (для обратной совместимости)
		if legacyID, ok := session.Data["task_id"].(string); ok {
			id, err := strconv.ParseInt(legacyID, 10, 64)
			if err != nil {
				return errors.New("Не удалось обновить задачу")
			}
			taskId = id
		} else {
			return errors.New("ID задачи не найден в сессии")
		}
	}

	// Получаем новое название задачи из сессии
	/* 	newTaskTitle, ok := session.Data["title"].(string)
	   	if !ok {
	   		return errors.New("arguments not found or invalid in session")
	   	} */
	// Получаем новое название задачи ИЗ СООБЩЕНИЯ (пользователь ввёл текст)
	newTaskTitle := msg.Text
	if newTaskTitle == "" {
		return errors.New("Новый текст задачи пустой")
	}

	// Парсим description и date
	taskDesc, taskDate, ok := messages.ParseForCommandSaveTask(newTaskTitle)
	if !ok {
		return errors.New("Не удалось сохранить задачу")
	}

	// Если описание задачи пустое
	if taskDesc == "" {
		return errors.New("Описание задачи отсутствует")
	}

	// Проверяем, что ID задачи является числом
	/* 	argsInt, err := strconv.Atoi(strings.TrimSpace(taskId))
	   	if err != nil {
	   		return errors.Wrap(err, "Не удалось удалить задачу")
	   	} */

	// Получаем задачу по ID
	// task, err = s.GetTaskStorage().GetTaskByID(s.GetCtx(), argsInt)
	task, err := s.GetTaskStorage().GetTaskByID(s.GetCtx(), taskId)
	if err != nil {
		return errors.Wrap(err, "Не удалось изменить статус задачи")
	}

	// Получаем дату обновления задачи
	timestamp := time.Unix(msg.Date, 0)

	// Обновляем description и updated_at задачи
	task.Description = taskDesc
	task.UpdatedAt = timestamp

	// Если дата не пустая обновляем deadline задачи
	if taskDate != "" {
		layout := "02.01.2006 15:04"
		task.Deadline, _ = time.Parse(layout, taskDate)
	}

	// Сохраняем задачу
	err = s.GetTaskStorage().UpdateTask(s.GetCtx(), task)
	if err != nil {
		return errors.Wrap(err, "Не удалось сохранить задачу")
	}

	// Удаляем сессию
	_ = s.GetSessionService().DeleteSession(s.GetCtx(), msg.UserID)
	//if err != nil {
	//	return errors.Wrap(err, "Не удалось сохранить задачу")
	//}

	return s.GetTgClient().ShowInlineButtons(resources.TXTUpdateTask, btnupdatetask.BtnUpdateTask, msg.UserID)
}
