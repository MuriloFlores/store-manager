package config

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
)

var Cfg *Config

type Config struct {
	DBHost                   string `mapstructure:"DB_HOST"`
	DBPort                   string `mapstructure:"DB_PORT"`
	DBUser                   string `mapstructure:"DB_USER"`
	DBPassword               string `mapstructure:"DB_PASSWORD"`
	DBName                   string `mapstructure:"DB_NAME"`
	ServerPort               string `mapstructure:"SERVER_PORT"`
	JWTSecret                string `mapstructure:"JWT_SECRET"`
	RedisAddress             string `mapstructure:"REDIS_ADDRESS"`
	SmtpHost                 string `mapstructure:"SMTP_HOST"`
	SmtpPort                 string `mapstructure:"SMTP_PORT"`
	SmtpSenderEmail          string `mapstructure:"SMTP_SENDER_EMAIL"`
	SmtpAppPassword          string `mapstructure:"SMTP_APP_PASSWORD"`
	FrontendPasswordResetURL string `mapstructure:"FRONTEND_PASSWORD_RESET_URL"`
}

func LoadConfig(path string) (*Config, error) {
	//viper.AddConfigPath(path)
	//viper.SetConfigName(".env")
	//viper.SetConfigType("env")

	//viper.ReadInConfig()

	viper.BindEnv("DB_HOST")
	viper.BindEnv("DB_PORT")
	viper.BindEnv("DB_USER")
	viper.BindEnv("DB_PASSWORD")
	viper.BindEnv("DB_NAME")
	viper.BindEnv("SERVER_PORT")
	viper.BindEnv("JWT_SECRET")
	viper.BindEnv("REDIS_ADDRESS")
	viper.BindEnv("SMTP_HOST")
	viper.BindEnv("SMTP_PORT")
	viper.BindEnv("SMTP_SENDER_EMAIL")
	viper.BindEnv("SMTP_APP_PASSWORD")
	viper.BindEnv("FRONTEND_PASSWORD_RESET_URL")

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("não foi possível fazer unmarshal da configuração: %w", err)
	}

	if cfg.DBUser == "" || cfg.DBPassword == "" || cfg.DBName == "" {
		return nil, errors.New("variáveis de banco de dados essenciais (DB_USER, DB_PASSWORD, DB_NAME) não foram carregadas")
	}

	Cfg = &cfg
	return Cfg, nil
}
func (c *Config) DatabaseURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=require",
		c.DBUser,
		c.DBPassword,
		c.DBHost,
		c.DBPort,
		c.DBName,
	)
}
