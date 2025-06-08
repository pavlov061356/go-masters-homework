package config

import (
	"errors"

	"github.com/spf13/viper"
)

var (
	ErrDBPathRequired = errors.New("db_path is required")
	ErrPortInvalid    = errors.New("port must be between 1 and 65535")
)

type Config struct {
	DBPath string `mapstructure:"db_path"`
	Port   int    `mapstructure:"port"`

	// AvgScoreRefreshTime -- время пересчёта средней оценки в секундах
	// Во время пересчёта средней оценки при добавлении новой оценки
	// накапливается ошибка, персчёт её сбрасывает
	AvgScoreRefreshTime int64 `mapstructure:"avg_score_refresh_time"`

	SentimenterQueue SentimenterQueue `mapstructure:"sentimenter_queue"`

	Sentimenter Sentimenter `mapstructure:"sentimenter"`
}

// SentimenterQueue - структура настроек очереди обработки настроения отзывов пользователей
type SentimenterQueue struct {
	// MaxSentimenterQueueLen - максимальная длина очереди ожидания записи в БД
	MaxDBQueueLen int `mapstructure:"max_db_queue_len"`
	// MaxDBQueueWait - максимальное время ожидания получения новых отзывов пользователей в секундах
	MaxDBQueueWait int `mapstructure:"max_db_queue_wait"`
}

// Sentimenter - структура настроек подключения к сервису расчёта настроения отзывов пользователей
type Sentimenter struct {
	// Addr - адрес сервиса расчёта настроения отзывов пользователей
	Addr string `mapstructure:"addr"`
	// Timeout - таймаут запроса к сервису расчёта настроения отзыва пользователя в секундах
	Timeout int `mapstructure:"timeout"`
	// Model - название модели нейросети для расчёта настроения отзыва пользователя
	Model string `mapstructure:"model"`
}

func (c *Config) validate() error {
	if c.DBPath == "" {
		return ErrDBPathRequired
	}

	if c.Port <= 0 || c.Port > 65535 {
		return ErrPortInvalid
	}

	return nil
}

func defaults() {
	viper.SetDefault("port", 8080)
	viper.SetDefault("avg_score_refresh_time", 604800) // 1 неделя в секундах
	viper.SetDefault("sentimenter_queue.max_db_queue_len", 1000)
	viper.SetDefault("sentimenter_queue.max_db_queue_wait", 600) // 10 минут
	viper.SetDefault("sentimenter.timeout", 60)                  // 1 минута
}

func Load(path string) (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	if path != "" {
		viper.SetConfigFile(path)
	}

	viper.AddConfigPath(".")
	defaults()

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = viper.Unmarshal(&cfg)

	if err != nil {
		return nil, err
	}

	err = cfg.validate()
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
