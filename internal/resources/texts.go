// Package resources
package resources

const (
	TXTStart                    = "<b>%v</b>. Выберите действие:"
	TXTViewTasks                = "<b>%v</b> выберите фильтр для задач:"
	TXTAddTask                  = "<b>%v</b>, введите описание задачи и, если нужно, дедлайн:\nКупить молоко 10.10.2026 18:00\nКупить молоко 10.10.2026 (до конца дня)\nКупить молоко (без дедлайна)"
	TXTAddTaskTemplate          = "Купить молоко 10.10.2026 18:00"
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
	ErrTaskDescriptionEmpty   = "Введите описание задачи, например:\nКупить молоко 10.10.2026 18:00"
	ErrArgumentsNotFound      = "Аргументы не найдены или некорректны"
	ErrFailedToGetUser        = "Не удалось получить пользователя"
	ErrFailedToUpdateTimezone = "Не удалось сохранить настройки. Попробуйте позже."
	ErrInvalidTimezone        = "Некорректный часовой пояс. Пожалуйста, попробуйте ещё раз."
	ErrSessionDataMissing     = "Данные сессии отсутствуют"
	ErrGeneric                = "Что-то пошло не так. Попробуйте ещё раз позже."
	ErrInvalidDeadline        = "Такой даты или времени не существует. Укажите дедлайн в формате:\nКупить молоко 10.10.2026 18:00"
	ErrTimeWithoutDate        = "Укажите дату вместе со временем, например:\nКупить молоко 10.10.2026 18:00"
	ErrDeadlinePassed         = "Дедлайн уже прошёл. Укажите дату в будущем."
)
