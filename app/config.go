package app

import (
	"fmt"
	"github.com/spf13/viper"
)

type serviceConfig struct {
	Sample          string
	Environment     string
	pyServerBaseURL string
}

func InitServiceConfig() serviceConfig {
	return serviceConfig{
		Sample:          ReadEnvString("SAMPLE"),
		Environment:     ReadEnvString("ENVIRONMENT"),
		pyServerBaseURL: ReadEnvString("PYTHON_SERVER_BASE_URL"),
	}
}

func (s *serviceConfig) GetSample() string {
	return s.Sample
}

func (s *serviceConfig) GetEnv() string {
	return s.Environment
}

func (s *serviceConfig) GetPythonServerBaseURL() string {
	return s.pyServerBaseURL
}

func ReadEnvString(key string) string {
	CheckIfSet(key)
	return viper.GetString(key)
}

func CheckIfSet(key string) {
	if !viper.IsSet(key) {
		err := fmt.Errorf("key %s is not set", key)
		panic(err)
	}
}
