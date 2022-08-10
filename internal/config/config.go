package config

import (
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
	"go.uber.org/zap"
)

type Config struct {
	IsDebug *bool `yaml:"is_debug"`
	Ports   struct {
		HttpPort string `yaml:"http_port"`
		GrpcPort string `yaml:"grpc_port"`
	}
	Hosts struct {
		AuthHost string `yaml:"auth_host"`
	}
}

var instance *Config
var once sync.Once

func GetConfig(logger *zap.SugaredLogger) *Config {
	once.Do(func() {
		instance = &Config{}
		if err := cleanenv.ReadConfig("config.yml", instance); err != nil {
			help, _ := cleanenv.GetDescription(instance, nil)
			logger.Errorf(help) // TODO: Use normal logger here !!!
		}
	})
	return instance
}
