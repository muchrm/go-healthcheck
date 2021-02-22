package lineauth

import (
	"testing"

	"github.com/muchrm/go-healthcheck/config"
	"github.com/spf13/viper"
)

func Test_bindLocalServer(t *testing.T) {
	t.Run("should be throw when not set AuthHostAddr", func(t *testing.T) {
		viper.Set(config.AuthHostAddr, nil)
		_, err := bindLocalServer()
		if err == nil {
			t.Errorf("bindLocalServer() error = %v, wantErr %v", err, true)
			return
		}
	})
	t.Run("should be throw when AuthHostAddr incorect", func(t *testing.T) {
		viper.Set(config.AuthHostAddr, "ssllldd.0.0.3")
		_, err := bindLocalServer()
		if err == nil {
			t.Errorf("bindLocalServer() error = %v, wantErr %v", err, true)
			return
		}
	})
	t.Run("should not be throw when AuthHostAddr corect", func(t *testing.T) {
		viper.Set(config.AuthHostAddr, "127.0.0.1:8081")
		server, err := bindLocalServer()
		server.Close()
		if err != nil {
			t.Errorf("InitFlow() error = %v, wantErr %v", err, false)
			return
		}
	})
}
