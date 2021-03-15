package config

import (
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	Dsn string
	Port string

}

func Load() *Config {
	conf := &Config{}
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error while reading config file %s", err)
	}
	conf.Dsn = viper.Get("DB_DSN").(string)
	conf.Port = viper.Get("PORT").(string)

	return conf
}