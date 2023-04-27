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
	Logger                LoggerConf    `mapstructure:"Logger"`
	ServerShutdownTimeout time.Duration `mapstructure:"SERVER_SHUTDOWN_TIMEOUT"`
	dbConnMaxLifetime     time.Duration `mapstructure:"DB_CONN_MAX_LIFETIME"`
	dbTimeOut             time.Duration `mapstructure:"DB_TIMEOUT"`
	address               string        `mapstructure:"ADDRESS"`
	port                  string        `mapstructure:"PORT"`
	grpcPort              string        `mapstructure:"GRPC_PORT"`
	OSFilePathSeparator   string
	dbAddress 			  string `mapstructure:"DB_ADDRESS"`
	dbPort				  string `mapstructure:"DB_PORT"`
	dbName                string `mapstructure:"MYSQL_DATABASE"`
	dbUser                string `mapstructure:"MYSQL_USER"`
	dbPassword            string `mapstructure:"MYSQL_PASSWORD"`
	dbMaxOpenConns        int    `mapstructure:"DB_MAX_OPEN_CONNS"`
	dbMaxIdleConns        int    `mapstructure:"DB_MAX_IDLE_CONNS"`
	workWithDBStorage     bool   `mapstructure:"WORK_W_DB_STORAGE"`
}

type LoggerConf struct {
	Level string `mapstructure:"LOG_LEVEL"`
}

func NewConfig() Config {
	return Config{}
}

func (config *Config) Init(path string) error {
	if path == "" {
		err := errors.New("void path to config.env")
		return err
	}

	viper.SetDefault("ADDRESS", "127.0.0.1")
	viper.SetDefault("PORT", "4000")
	viper.SetDefault("GRPC_PORT", "5000")
	viper.SetDefault("SERVER_SHUTDOWN_TIMEOUT", 30*time.Second)
	viper.SetDefault("MYSQL_DATABASE", "OTUSFinalLab")
	viper.SetDefault("MYSQL_USER", "imapp")
	viper.SetDefault("MYSQL_PASSWORD", "LightInDark")
	viper.SetDefault("DB_CONN_MAX_LIFETIME", time.Minute*3)
	viper.SetDefault("DB_MAX_OPEN_CONNS", 20)
	viper.SetDefault("DB_MAX_IDLE_CONNS", 20)
	viper.SetDefault("DB_TIMEOUT", 5*time.Second)
	viper.SetDefault("WORK_W_DB_STORAGE", true)
	viper.SetDefault("DB_ADDRESS ", "127.0.0.1")
	viper.SetDefault("DB_PORT", "3306")
	viper.SetDefault("LOG_LEVEL", "debug")

	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok { //nolint:errorlint
			return err
		}
	}

	viper.Set("OSFilePathSeparator", string(os.PathSeparator))

	config.address = viper.GetString("ADDRESS")
	config.port = viper.GetString("PORT")
	config.grpcPort = viper.GetString("GRPC_PORT")
	config.ServerShutdownTimeout = viper.GetDuration("SERVER_SHUTDOWN_TIMEOUT")
	config.dbName = viper.GetString("MYSQL_DATABASE")
	config.dbUser = viper.GetString("MYSQL_USER")
	config.dbPassword = viper.GetString("MYSQL_PASSWORD")
	config.dbConnMaxLifetime = viper.GetDuration("DB_CONN_MAX_LIFETIME")
	config.dbMaxOpenConns = viper.GetInt("DB_MAX_OPEN_CONNS")
	config.dbMaxIdleConns = viper.GetInt("DB_MAX_IDLE_CONNS")
	config.dbTimeOut = viper.GetDuration("DB_TIMEOUT")
	config.workWithDBStorage = viper.GetBool("WORK_W_DB_STORAGE")
	config.Logger.Level = viper.GetString("LOG_LEVEL")
	config.dbAddress = viper.GetString("DB_ADDRESS")
	config.dbPort = viper.GetString("DB_PORT")
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

func (config *Config) GetGRPCPort() string {
	return config.grpcPort
}

func (config *Config) GetOSFilePathSeparator() string {
	return config.OSFilePathSeparator
}

func (config *Config) GetServerShutdownTimeout() time.Duration {
	return config.ServerShutdownTimeout
}

func (config *Config) GetDBName() string {
	return config.dbName
}

func (config *Config) GetDBUser() string {
	return config.dbUser
}

func (config *Config) GetDBPassword() string {
	return config.dbPassword
}

func (config *Config) GetDBConnMaxLifetime() time.Duration {
	return config.dbConnMaxLifetime
}

func (config *Config) GetDBMaxOpenConns() int {
	return config.dbMaxOpenConns
}

func (config *Config) GetDBMaxIdleConns() int {
	return config.dbMaxIdleConns
}

func (config *Config) GetDBTimeOut() time.Duration {
	return config.dbTimeOut
}

func (config *Config) GetDBAddress() string {
	return config.dbAddress
}

func (config *Config) GetDBPort() string {
	return config.dbPort
}