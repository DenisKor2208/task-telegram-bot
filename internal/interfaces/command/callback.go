package command

type CallbackHandler interface {
	Command

	//CallbackAction - Вернуть действие (например, "delete_task_by_id")
	CallbackAction() string
}
