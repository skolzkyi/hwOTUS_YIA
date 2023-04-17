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
	Logger                LoggerConf
	kafkaTimeOut          time.Duration `mapstructure:"S_KAFKA_TIMEOUT"`
	senderShutdownTimeout time.Duration `mapstructure:"SENDER_SHUTDOWN_TIMEOUT"`
	kafkaAddr             string        `mapstructure:"S_KAFKA_ADDR"`
	kafkaPort             string        `mapstructure:"S_KAFKA_PORT"`
	kafkaTopicName        string        `mapstructure:"KAFKA_CREATE_TOPICS"`
}

type LoggerConf struct {
	Level string `mapstructure:"LOG_LEVEL"`
}

func NewConfig() Config {
	return Config{}
}

func (config *Config) Init(path string) error {
	if path == "" {
		err := errors.New("void path to config_sender.env")
		return err
	}

	viper.SetDefault("S_KAFKA_ADDR", "127.0.0.1")
	viper.SetDefault("S_KAFKA_PORT", "9092")
	viper.SetDefault("KAFKA_CREATE_TOPICS", "CLNotifications1")
	viper.SetDefault("S_KAFKA_TIMEOUT", 3*time.Second)
	viper.SetDefault("SENDER_SHUTDOWN_TIMEOUT", 30*time.Second)

	viper.SetDefault("LOG_LEVEL", "debug")

	viper.AddConfigPath(path)
	viper.SetConfigName("config_sender")
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
	config.senderShutdownTimeout = viper.GetDuration("SENDER_SHUTDOWN_TIMEOUT")

	config.Logger.Level = viper.GetString("LOG_LEVEL")

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

func (config *Config) GetSenderShutdownTimeout() time.Duration {

	return config.senderShutdownTimeout

}

func (config *Config) GetKafkaTimeOut() time.Duration {

	return config.kafkaTimeOut
}
