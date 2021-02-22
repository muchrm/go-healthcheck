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
	AccessToken                string = "AccessToken"
	HealcheckReportAPI         string = "CSVPathMissing"
	TypeMissMatchError         string = "TypeMissMatch"
	CSVNotFoundError           string = "FileNotFound"
	CSVIncorrectFormatError    string = "CSVIncorrectFormat"
	CSVPathMissingError        string = "CSVPathMissing"
	AuthClientID               string = "ClientId"
	AuthClientSecret           string = "ClientSecret"
	AuthRedirectURI            string = "RedirectURI"
	AuthHostAddr               string = "AuthHostAddr"
	AuthHost                   string = "AuthHost"
	AuthResourceHose           string = "AuthResourceHose"
	AuthScopes                 string = "Scopes"
)

func InitConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath("$HOME/.config/go-healthcheck")
	viper.AddConfigPath(".")
	viper.SetEnvPrefix("healthcheck")
	viper.BindEnv(AccessToken)

	args := os.Args[1:]
	if len(args) > 0 && args[0] != "" {
		viper.SetDefault(CSVPath, args[0])
	}

	viper.SetDefault(MaxWorker, runtime.NumCPU())
	viper.SetDefault(HealcheckReportAPI, "https://backend-challenge.line-apps.com/healthcheck/report")
	viper.SetDefault(AuthHost, "https://access.line.me/oauth2/v2.1")
	viper.SetDefault(AuthResourceHose, "https://api.line.me/oauth2/v2.1")
	viper.SetDefault(AuthScopes, []string{"openid"})
	viper.SetDefault(AuthRedirectURI, "http://127.0.0.1:8081/auth")
	viper.SetDefault(AuthHostAddr, "127.0.0.1:8081")
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

func GetStringArray(configKey string) ([]string, error) {
	configValue, ok := viper.Get(configKey).([]string)
	if !ok {
		return []string{}, errors.New(TypeMissMatchError)
	}
	return configValue, nil
}
