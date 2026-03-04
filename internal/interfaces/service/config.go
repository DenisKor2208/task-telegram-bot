package service

type ConfigEntry struct {
	TargetCommand string
	CommandFields []string
}

type ConfigEntryService interface {
	GetConfig() map[string]ConfigEntry
}
