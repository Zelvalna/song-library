package config

import "github.com/spf13/viper"

type Config struct {
	Port   string
	DB     DBConfig
	APIURL string
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile(".env")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	config := &Config{
		Port:   viper.GetString("PORT"),
		APIURL: viper.GetString("API_URL"),
		DB: DBConfig{
			Host:     viper.GetString("DB_HOST"),
			Port:     viper.GetString("DB_PORT"),
			User:     viper.GetString("DB_USER"),
			Password: viper.GetString("DB_PASSWORD"),
			Name:     viper.GetString("DB_NAME"),
		},
	}

	return config, nil
}
