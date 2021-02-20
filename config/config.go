package config

import (
	"errors"
	"os"
	"runtime"

	"github.com/spf13/viper"
)

var (
	MaxWorker                  string = "MaxWorker"
	CSVPath                    string = "CSVPath"
	HealthcheckTimeoutInSecond string = "HealthcheckTimeoutInSecond"
	LineAPIToken               string = "LineAPIToken"
	HealcheckReportAPI         string = "CSVPathMissing"
	TypeMissMatchError         string = "TypeMissMatch"
	CSVNotFoundError           string = "FileNotFound"
	CSVIncorrectFormatError    string = "CSVIncorrectFormat"
	CSVPathMissingError        string = "CSVPathMissing"
)

func InitConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath("$HOME/.config/go-healthcheck")
	viper.AddConfigPath(".")
	viper.SetEnvPrefix("healthcheck")
	viper.BindEnv(LineAPIToken)

	args := os.Args[1:]
	if len(args) > 0 && args[0] != "" {
		viper.SetDefault(CSVPath, args[0])
	}

	viper.SetDefault(MaxWorker, runtime.NumCPU())
	viper.SetDefault(HealcheckReportAPI, "https://backend-challenge.line-apps.com/healthcheck/report")
	viper.ReadInConfig()

}

func GetInt(configKey string) (int, error) {
	configValue, ok := viper.Get(configKey).(int)
	if !ok {
		return 0, errors.New(TypeMissMatchError)
	}
	return configValue, nil
}
func GetString(configKey string) (string, error) {
	configValue, ok := viper.Get(configKey).(string)
	if !ok {
		return "", errors.New(TypeMissMatchError)
	}
	return configValue, nil
}
