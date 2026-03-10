package config

import (
	"flag"
	"os"

	"github.com/DenisKor2208/task-telegram-bot/internal/logger"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/pkg/errors"
)

const configFile = "./config/local.yaml"

type RedisConfig struct {
	ConnectionString string `yaml:"REDIS_CONNECTION_STRING" env:"REDIS_CONNECTION_STRING" env-default:"redis://localhost:6379"`
	SessionTTL       int    `yaml:"REDIS_SESSION_TTL" env:"REDIS_SESSION_TTL" env-default:"3600"` // Время жизни сессии в секундах (3600 = 1 час)
}

type Config struct {
	Env                  string      `yaml:"env" env-default:"local"` // Текущий окружение
	APIBotToken          string      `yaml:"api_telegram_token"`      // Токен бота в телеграме
	ConnectionStringDB   string      `yaml:"ConnectionStringDB"`      // Строка подключения в базе данных.
	Redis                RedisConfig `yaml:"redis"`
	OverdueCheckInterval int         `yaml:"overdue_check_interval" env-default:"60"` // в секундах
}

type Service struct {
	config Config
}

func New() (*Service, error) {
	path := fetchConfigPath()

	s := &Service{}

	if _, err := os.Stat(path); err != nil {
		logger.Error("Ошибка init config file", "err", err)
		return nil, errors.Wrap(err, "init config file")
	}

	if err := cleanenv.ReadConfig(path, &s.config); err != nil {
		logger.Error("Ошибка при чтении config-файла", "err", err)
		return nil, errors.Wrap(err, "parsing yaml")
	}

	return s, nil
}

// fetchConfigPath - получаем путь к конфигу
func fetchConfigPath() string {
	var res string

	// Получаем путь к конфигу из флага
	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	// Получаем путь к конфигу из переменной окружения
	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	// Получаем путь к конфигу из константы
	if res == "" {
		res = configFile
	}

	return res
}

// Token - получаем токен бота в Телеграм
func (s *Service) Token() string {
	return s.config.APIBotToken
}

// GetConfig - получаем конфиг
func (s *Service) GetConfig() Config {
	return s.config
}
