package app

import (
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"taleteller/db"
	"taleteller/logger"
)

var context *appContext

type appContext struct {
	zapLogger     *zap.SugaredLogger
	ServiceConfig serviceConfig
}

func Init() (err error) {
	context = &appContext{}

	err = initConfig()
	if err != nil {
		return
	}

	context.ServiceConfig = InitServiceConfig()
	fmt.Println("yml - ", context.ServiceConfig.GetSample())

	zapLogger, err := logger.SetupLogger(context.ServiceConfig.GetEnv())
	if err != nil {
		return
	}

	context.zapLogger = zapLogger

	err = initDB()
	if err != nil {
		return
	}

	return
}

func Close() (err error) {
	// flushes log buffer, if any
	if context.zapLogger != nil {
		context.zapLogger.Sync()
	}
	return
}

func initConfig() (err error) {

	err = Load("./", "application")
	if err != nil {
		return
	}

	return
}
func Load(filePath string, fileName string) (err error) {
	viper.SetDefault("ENV", "development")
	viper.SetConfigName(fileName)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(filePath)
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if _, ok := err.(viper.ConfigFileNotFoundError); err != nil && !ok {
		return
	}

	//appConfig = config{}

	return
}

func initDB() (err error) {
	user := "postgres"
	password := "postgres"
	host := "localhost"
	port := 5432
	dbName := "dpay"
	err = db.Init(&db.Config{
		Driver:       "postgres",
		URL:          ConnectionURL(user, password, host, port, dbName),
		MaxIdleConns: 100,
		MaxOpenConns: 5,
	})
	if err != nil {
		return
	}

	return
}

func ConnectionURL(user string, password string, host string, port int, dbName string) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", user, password, host, port, dbName)
}
