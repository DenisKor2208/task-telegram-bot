// Package resources
package resources

const (
	TXTStart                    = "<b>%v</b>. Выберите действие:"
	TXTViewTasks                = "<b>%v</b> выберите фильтр для задач:"
	TXTAddTask                  = "<b>%v</b> введите описание задачи и дедлайн в формате:\nВыполнить 1 пункт 01.01.2025 01:00"
	TXTAddTaskTemplate          = "Выполнить 1 пункт ТЗ 01.01.2025 01:00"
	TXTSaveTask                 = "Задача успешно добавлена"
	TXTFilterAll                = "<b>%v</b> вы просматриваете все имеющиеся задачи:"
	TXTFilterCompleted          = "<b>%v</b> вы просматриваете все имеющиеся задачи в статусе <b>Выполнено</b>:"
	TXTFilterOverdue            = "<b>%v</b> вы просматриваете все имеющиеся задачи в статусе <b>Просрочено</b>:"
	TXTFilterClosed             = "<b>%v</b> вы просматриваете все имеющиеся задачи в статусе <b>Завершено</b>:"
	TXTFilterInProgress         = "<b>%v</b> вы просматриваете все имеющиеся задачи в статусе <b>В процессе</b>:"
	TXTDeleteTask               = "<b>%v</b> выберите задачу для удаления:"
	TXTCompletedTaskByIDCommand = "<b>%v</b> здесь нужно добавить соответствующий текст:"
	TXTClosedTaskByIDCommand    = "<b>%v</b> здесь нужно добавить соответствующий текст:"
	TXTDeleteTaskByIDCommand    = "<b>%v</b> здесь нужно добавить соответствующий текст:"
	TXTSetTimezone              = "<b>%v</b> пожалуйста, выберите ваш часовой пояс:"
	TXTSetTimezoneConfirm       = "✅ Часовой пояс успешно сохранён.\nТеперь даты будут интерпретироваться в вашем локальном времени."
	TXTCompletedTask            = "<b>%v</b> выберите задачу для пометки <b>Выполнено</b>:"
	TXTClosedTask               = "<b>%v</b> выберите задачу для пометки <b>Завершено</b>:"
	TXTUpdateTask               = "Задача успешно обновлена"
	TXTUnknownCommand           = "К сожалению, данная команда мне неизвестна. Для начала работы введите /start"
)

// FilterTexts сопоставляет ключи команд фильтрации с соответствующими текстовыми константами.
var FilterTexts = map[string]string{
	"filter_all":         TXTFilterAll,
	"filter_in_progress": TXTFilterInProgress,
	"filter_completed":   TXTFilterCompleted,
	"filter_overdue":     TXTFilterOverdue,
	"filter_closed":      TXTFilterClosed,
}

// Ошибки
const (
	ErrInvalidTaskID          = "Неверный ID задачи"
	ErrTaskNotFound           = "Задача не найдена"
	ErrFailedToUpdateTask     = "Не удалось изменить статус задачи"
	ErrFailedToDeleteTask     = "Не удалось удалить задачу"
	ErrFailedToSaveTask       = "Не удалось сохранить задачу"
	ErrTaskDescriptionEmpty   = "Описание задачи отсутствует"
	ErrArgumentsNotFound      = "Аргументы не найдены или некорректны"
	ErrFailedToGetUser        = "Не удалось получить пользователя"
	ErrFailedToUpdateTimezone = "Не удалось сохранить настройки. Попробуйте позже."
	ErrInvalidTimezone        = "Некорректный часовой пояс. Пожалуйста, попробуйте ещё раз."
	ErrSessionDataMissing     = "Данные сессии отсутствуют"
)
