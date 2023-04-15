package main

import (
	"errors"
	"time"

	"github.com/spf13/viper"
)

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type Config struct {
	Logger                     LoggerConf
	kafkaTimeOut               time.Duration `mapstructure:"S_KAFKA_TIMEOUT"`
	shedulerShutdownTimeout    time.Duration `mapstructure:"SHEDULER_SHUTDOWN_TIMEOUT"`
	notificationEventPeriod    time.Duration `mapstructure:"NOTIFICATION_EVENT_PERIOD"`
	cleanOldEventPeriod        time.Duration `mapstructure:"CLEAN_OLD_EVENT_PERIOD"`
	shedulerPeriod             time.Duration `mapstructure:"SHEDULER_PERIOD"`
	kafkaAddr                  string        `mapstructure:"S_KAFKA_ADDR"`
	kafkaPort                  string        `mapstructure:"S_KAFKA_PORT"`
	kafkaTopicName             string        `mapstructure:"KAFKA_CREATE_TOPICS"`
	grpsport                   string        `mapstructure:"GRPC_PORT"`
	serverURL                  string        `mapstructure:"SERVER_URL"`
	kafkaAutoCreateTopicEnable bool          `mapstructure:"KAFKA_AUTO_CREATE_TOPICS_ENABLE"`
}

type LoggerConf struct {
	Level string `mapstructure:"LOG_LEVEL"`
}

func NewConfig() Config {
	return Config{}
}

func (config *Config) Init(path string) error {
	if path == "" {
		err := errors.New("void path to config_sheduler.env")
		return err
	}

	viper.SetDefault("S_KAFKA_ADDR", "127.0.0.1")
	viper.SetDefault("S_KAFKA_PORT", "9092")
	viper.SetDefault("KAFKA_CREATE_TOPICS", "CLNotification1")
	viper.SetDefault("S_KAFKA_TIMEOUT", 3*time.Second)
	viper.SetDefault("SHEDULER_SHUTDOWN_TIMEOUT", 30*time.Second)
	viper.SetDefault("NOTIFICATION_EVENT_PERIOD", 30*time.Second)
	viper.SetDefault("CLEAN_OLD_EVENT_PERIOD", 30*time.Second)
	viper.SetDefault("SHEDULER_PERIOD", 10*time.Second)
	viper.SetDefault("GRPC_PORT", "5000")
	viper.SetDefault("SERVER_URL", "127.0.0.1")

	viper.SetDefault("LOG_LEVEL", "debug")
	viper.SetDefault("KAFKA_AUTO_CREATE_TOPICS_ENABLE", true)

	viper.AddConfigPath(path)
	viper.SetConfigName("config_sheduler")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}

	config.kafkaAddr = viper.GetString("S_KAFKA_ADDR")
	config.kafkaPort = viper.GetString("S_KAFKA_PORT")
	config.kafkaTopicName = viper.GetString("KAFKA_CREATE_TOPICS")
	config.kafkaTimeOut = viper.GetDuration("S_KAFKA_TIMEOUT")
	config.shedulerShutdownTimeout = viper.GetDuration("SHEDULER_SHUTDOWN_TIMEOUT")
	config.notificationEventPeriod = viper.GetDuration("NOTIFICATION_EVENT_PERIOD")
	config.cleanOldEventPeriod = viper.GetDuration("CLEAN_OLD_EVENT_PERIOD")
	config.shedulerPeriod = viper.GetDuration("SHEDULER_PERIOD")
	config.Logger.Level = viper.GetString("LOG_LEVEL")
	config.grpsport = viper.GetString("GRPC_PORT")
	config.serverURL = viper.GetString("SERVER_URL")
	config.kafkaAutoCreateTopicEnable = viper.GetBool("KAFKA_AUTO_CREATE_TOPICS_ENABLE")

	return nil

}

func (config *Config) GetKafkaURL() string {

	return config.kafkaAddr + ":" + config.kafkaPort

}

func (config *Config) GetKafkaAddr() string {

	return config.kafkaAddr

}

func (config *Config) GetKafkaPort() string {

	return config.kafkaPort

}

func (config *Config) GetKafkaTopicName() string {

	return config.kafkaTopicName

}

func (config *Config) GetShedulerShutdownTimeout() time.Duration {

	return config.shedulerShutdownTimeout

}

func (config *Config) GetNotificationEventPeriod() time.Duration {

	return config.notificationEventPeriod
}

func (config *Config) GetKafkaTimeOut() time.Duration {

	return config.kafkaTimeOut
}

func (config *Config) GetCleanOldEventPeriod() time.Duration {

	return config.cleanOldEventPeriod
}

func (config *Config) GetShedulerPeriod() time.Duration {

	return config.shedulerPeriod
}

func (config *Config) GetGRPSPort() string {

	return config.grpsport

}

func (config *Config) GetServerURL() string {

	return config.serverURL

}

func (config *Config) GetKafkaAutoCreateTopicEnable() bool {

	return config.kafkaAutoCreateTopicEnable
}
