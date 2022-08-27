package app

import (
	"fmt"
	"github.com/spf13/viper"
)

type serviceConfig struct {
	Sample                   string
	Environment              string
	pyServerBaseURL          string
	awsRegion                string
	awsAccessKeyID           string
	awsSecretAccessKey       string
	awsGeneratedAssetsBucket string
}

func InitServiceConfig() serviceConfig {
	return serviceConfig{
		Sample:                   ReadEnvString("SAMPLE"),
		Environment:              ReadEnvString("ENVIRONMENT"),
		pyServerBaseURL:          ReadEnvString("PYTHON_SERVER_BASE_URL"),
		awsRegion:                ReadEnvString("AWS_REGION"),
		awsAccessKeyID:           ReadEnvString("AWS_ACCESS_KEY_ID"),
		awsSecretAccessKey:       ReadEnvString("AWS_SECRET_ACCESS_KEY"),
		awsGeneratedAssetsBucket: ReadEnvString("AWS_GENERATED_ASSETS_BUCKET"),
	}
}

func (s *serviceConfig) GetAWSRegion() string {
	return s.awsRegion
}

func (s *serviceConfig) GetAWSAccessKeyID() string {
	return s.awsAccessKeyID
}

func (s *serviceConfig) GetAWSGeneratedAssetsBucket() string {
	return s.awsGeneratedAssetsBucket
}

func (s *serviceConfig) GetAWSSecretKey() string {
	return s.awsSecretAccessKey
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
