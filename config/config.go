package config

import (
	"log"

	"github.com/spf13/viper"
)

type (
	Connection struct {
		Alias    string `mapstructure:"alias"`
		Default  bool   `mapstructure:"default"`
		Host     string `mapstructure:"host"`
		Port     uint64 `mapstructure:"port"`
		Name     string `mapstructure:"name"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		Driver   string `mapstructure:"driver"`
	}

	Schema struct {
		Directory string `mapstructure:"dir"`
	}

	Configuration struct {
		Connections []Connection `mapstructure:"connections"`
		Schema      Schema       `mapstructure:"schema"`
	}
)

func Config() Configuration {

	var v = viper.New()
	v.SetConfigName("pebble")
	v.SetConfigType("yml")
	v.AddConfigPath(".")
	err := v.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
	var configuration Configuration
	v.Unmarshal(&configuration)

	return configuration
}
