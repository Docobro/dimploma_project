package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type SQLConfig struct {
	Collation         string
	Password          string
	DBName            string
	User              string
	Host              string
	Timezone          string
	Port              string
	Charset           string
	ReadTimeout       time.Duration
	MaxOpenConns      int
	MaxIdleConns      int
	Timeout           time.Duration
	WriteTimeout      time.Duration
	ConnMaxLifetime   time.Duration
	InterpolateParams bool
	ParseTime         bool
}

type HTTPConfig struct {
	Url          string
	Host         string
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type App struct {
	Env     string
	Name    string
	Version string
}

type CryptoConfig struct {
	Url string
	Key string
}
type Config struct {
	App
	HTTPConfig
	SQLConfig
	CryptoConfig
}

func New(filepath string) (*Config, error) {
	viper.SetConfigFile(filepath)
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("config - New - viper.ReadInConfig: %w", err)
	}

	viper.SetDefault("app.name", "testapp")
	viper.SetDefault("app.env", "dev")
	viper.SetDefault("app.version", "v1")
	viper.SetDefault("http.read_timeout", "1s")
	viper.SetDefault("http.write_timeout", "1s")

	config := &Config{
		App: App{
			Name:    viper.GetString("app.name"),
			Env:     viper.GetString("app.env"),
			Version: viper.GetString("app.version"),
		},

		HTTPConfig: HTTPConfig{
			Url:          viper.GetString("http.url"),
			Host:         viper.GetString("http.host"),
			Port:         viper.GetString("http.port"),
			ReadTimeout:  viper.GetDuration("http.read_timeout"),
			WriteTimeout: viper.GetDuration("http.write_timeout"),
		},

		SQLConfig: SQLConfig{
			Host:              viper.GetString("sql.host"),
			Port:              viper.GetString("sql.port"),
			DBName:            viper.GetString("sql.database"),
			User:              viper.GetString("sql.user"),
			Password:          viper.GetString("sql.password"),
			ReadTimeout:       viper.GetDuration("sql.read_timeout"),
			WriteTimeout:      viper.GetDuration("sql.write_timeout"),
			Timeout:           viper.GetDuration("sql.timeout"),
			InterpolateParams: false,
			Charset:           "UTF-8",
			ParseTime:         true,
			Timezone:          "Europe/Moscow",
			Collation:         "",
		},
		CryptoConfig: CryptoConfig{
			Url: viper.GetString("crypto.url"),
			Key: viper.GetString("crypto.key"),
		},
	}

	return config, nil
}
