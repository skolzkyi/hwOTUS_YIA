package main

import (
	"errors"
	"os"
	"time"

	"github.com/spf13/viper"
)

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type Config struct {
	Logger                LoggerConf
	ServerShutdownTimeout time.Duration
	dbConnMaxLifetime     time.Duration
	dbTimeOut             time.Duration
	address               string
	port                  string
	OSFilePathSeparator   string
	dbName                string
	dbUser                string
	dbPassword            string
	dbMaxOpenConns        int
	dbMaxIdleConns        int
	workWithDBStorage     bool
}

type LoggerConf struct {
	Level string
}

func NewConfig() Config {
	return Config{}
}

func (config *Config) Init(path string) error {
	if path == "" {
		err := errors.New("void path to config.yaml")
		return err
	}

	viper.SetDefault("address", "127.0.0.1")
	viper.SetDefault("port", "4000")
	viper.SetDefault("ServerShutdownTimeout", 30*time.Second)
	viper.SetDefault("dbName", "OTUSFinalLab")
	viper.SetDefault("dbUser", "imapp")
	viper.SetDefault("dbPassword", "LightInDark")
	viper.SetDefault("dbConnMaxLifetime", time.Minute*3)
	viper.SetDefault("dbMaxOpenConns", 20)
	viper.SetDefault("dbMaxIdleConns", 20)
	viper.SetDefault("dbTimeOut", 5*time.Second)
	viper.SetDefault("workWithDBStorage", false)
	viper.SetDefault("logLevel", "debug")

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)

	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}

	viper.Set("OSFilePathSeparator", string(os.PathSeparator))

	config.address = viper.GetString("address")
	config.port = viper.GetString("port")
	config.ServerShutdownTimeout = viper.GetDuration("ServerShutdownTimeout")
	config.dbName = viper.GetString("dbName")
	config.dbUser = viper.GetString("dbUser")
	config.dbPassword = viper.GetString("dbPassword")
	config.dbConnMaxLifetime = viper.GetDuration("dbConnMaxLifetime")
	config.dbMaxOpenConns = viper.GetInt("dbMaxOpenConns")
	config.dbMaxIdleConns = viper.GetInt("dbMaxIdleConns")
	config.dbTimeOut = viper.GetDuration("dbTimeOut")
	config.workWithDBStorage = viper.GetBool("workWithDBStorage")
	config.Logger.Level = viper.GetString("logLevel")
	config.OSFilePathSeparator = viper.GetString("OSFilePathSeparator")

	return nil

}

func (config *Config) GetServerURL() string {

	return config.address + ":" + config.port

}

func (config *Config) GetAddress() string {

	return config.address

}

func (config *Config) GetPort() string {

	return config.port

}

func (config *Config) GetOSFilePathSeparator() string {

	return config.OSFilePathSeparator

}

func (config *Config) GetServerShutdownTimeout() time.Duration {

	return config.ServerShutdownTimeout
}

func (config *Config) GetDbName() string {

	return config.dbName

}
func (config *Config) GetDbUser() string {

	return config.dbUser

}
func (config *Config) GetDbPassword() string {

	return config.dbPassword

}

func (config *Config) GetdbConnMaxLifetime() time.Duration {

	return config.dbConnMaxLifetime
}

func (config *Config) GetDbMaxOpenConns() int {

	return config.dbMaxOpenConns

}

func (config *Config) GetDbMaxIdleConns() int {

	return config.dbMaxIdleConns

}

func (config *Config) GetdbTimeOut() time.Duration {

	return config.dbTimeOut
}
