package configs

import (
	"sync"
	"time"

	"log"

	"github.com/spf13/viper"
)

type Config struct {
	APP struct {
		HOST string `mapstructure:"HOST"`
		PORT string `mapstructure:"PORT"`
	} `mapstructure:"APP"`

	Cache struct {
		Host          string        `mapstructure:"HOST"`
		Port          string        `mapstructure:"PORT"`
		DialTimeout   time.Duration `mapstructure:"DIAL_TIMEOUT"`
		ReadTimeout   time.Duration `mapstructure:"READ_TIMEOUT"`
		WriteTimeout  time.Duration `mapstructure:"WRITE_TIMEOUT"`
		IdleTimeout   time.Duration `mapstructure:"IDLE_TIMEOUT"`
		MaxConnAge    time.Duration `mapstructure:"MAX_CONN_AGE"`
		MinIdleConns  int           `mapstructure:"MIN_IDLE_CONNS"`
		Namespace     int           `mapstructure:"NAMESPACE"`
		Password      string        `mapstructure:"PASSWORD"`
		CacheDuration time.Duration `mapstructure:"CACHE_DURATION"`
	} `mapstructure:"CACHE"`

	JWT struct {
		Secret string `mapstructure:"SECRET"`
	} `mapstructure:"JWT"`

	DB struct {
		MySQL struct {
			Read struct {
				Host                 string        `mapstructure:"HOST"`
				Port                 string        `mapstructure:"PORT"`
				Name                 string        `mapstructure:"NAME"`
				Username             string        `mapstructure:"USERNAME"`
				Password             string        `mapstructure:"PASSWORD"`
				ConnOpenMax          int           `mapstructure:"CONN_OPEN_MAX"`
				ConnIdleMax          int           `mapstructure:"CONN_IDLE_MAX"`
				ConnLifetimeMax      time.Duration `mapstructure:"CONN_LIFETIME_MAX"`
				AdditionalParameters string        `mapstructure:"ADDITIONAL_PARAMETERS"`
			} `mapstructure:"READ"`
			Write struct {
				Host                 string        `mapstructure:"HOST"`
				Port                 string        `mapstructure:"PORT"`
				Name                 string        `mapstructure:"NAME"`
				Username             string        `mapstructure:"USERNAME"`
				Password             string        `mapstructure:"PASSWORD"`
				ConnOpenMax          int           `mapstructure:"CONN_OPEN_MAX"`
				ConnIdleMax          int           `mapstructure:"CONN_IDLE_MAX"`
				ConnLifetimeMax      time.Duration `mapstructure:"CONN_LIFETIME_MAX"`
				AdditionalParameters string        `mapstructure:"ADDITIONAL_PARAMETERS"`
			} `mapstructure:"WRITE"`
		} `mapstructure:"MYSQL"`
	} `mapstructure:"DB"`
}

var (
	conf Config
	once sync.Once
)

// Get are responsible to load env and get data an return the struct
func Get() *Config {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()

	if err != nil {
		log.Fatal(err)
	}

	once.Do(func() {
		log.Println("Service configuration initialized.")
		err = viper.Unmarshal(&conf)
		if err != nil {
			log.Fatal(err)
		}
	})

	return &conf
}
