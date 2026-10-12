// Package resources
package resources

const (
	TXTStart                    = "<b>%v</b>, выберите действие:"
	TXTViewTasks                = "<b>%v</b>, выберите фильтр для задач:"
	TXTAddTask                  = "<b>%v</b>, введите описание задачи и, если нужно, дедлайн:\nКупить молоко 10.10.2026 18:00\nКупить молоко 10.10.2026 (до конца дня)\nКупить молоко (без дедлайна)"
	TXTAddTaskTemplate          = "Купить молоко 10.10.2026 18:00"
	TXTSaveTask                 = "Задача успешно добавлена"
	TXTFilterAll                = "<b>%v</b>, все ваши задачи:"
	TXTFilterCompleted          = "<b>%v</b>, задачи в статусе <b>Выполнено</b>:"
	TXTFilterOverdue            = "<b>%v</b>, задачи в статусе <b>Просрочено</b>:"
	TXTFilterClosed             = "<b>%v</b>, задачи в статусе <b>Завершено</b>:"
	TXTFilterInProgress         = "<b>%v</b>, задачи в статусе <b>В процессе</b>:"
	TXTDeleteTask               = "<b>%v</b>, выберите задачу для удаления:"
	TXTCompletedTaskByIDCommand = "✅ Задача «%s» отмечена как <b>Выполнено</b>."
	TXTClosedTaskByIDCommand    = "🏁 Задача «%s» отмечена как <b>Завершено</b>."
	TXTDeleteTaskByIDCommand    = "🗑 Задача «%s» удалена."
	TXTSetTimezone              = "<b>%v</b>, пожалуйста, выберите ваш часовой пояс:"
	TXTSetTimezoneConfirm       = "✅ Часовой пояс успешно сохранён.\nТеперь даты будут интерпретироваться в вашем локальном времени."
	TXTCompletedTask            = "<b>%v</b>, выберите задачу для пометки <b>Выполнено</b>:"
	TXTClosedTask               = "<b>%v</b>, выберите задачу для пометки <b>Завершено</b>:"
	TXTUnknownCommand           = "К сожалению, данная команда мне неизвестна. Для начала работы введите /start"
)

// Тексты пустых списков задач по фильтрам
const (
	TXTFilterAllEmpty        = "<b>%v</b>, у вас пока нет задач."
	TXTFilterCompletedEmpty  = "<b>%v</b>, задач в статусе <b>Выполнено</b> нет."
	TXTFilterOverdueEmpty    = "<b>%v</b>, задач в статусе <b>Просрочено</b> нет."
	TXTFilterClosedEmpty     = "<b>%v</b>, задач в статусе <b>Завершено</b> нет."
	TXTFilterInProgressEmpty = "<b>%v</b>, задач в статусе <b>В процессе</b> нет."
)

// Тексты пустых списков задач для действий
const (
	TXTCompletedTaskEmpty = "<b>%v</b>, нет задач, которые можно отметить <b>Выполнено</b>."
	TXTClosedTaskEmpty    = "<b>%v</b>, нет задач, которые можно отметить <b>Завершено</b>."
	TXTDeleteTaskEmpty    = "<b>%v</b>, у вас нет задач для удаления."
)

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
