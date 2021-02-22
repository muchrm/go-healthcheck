package lineauth

import (
	"reflect"
	"testing"

	"github.com/muchrm/go-healthcheck/config"
	"github.com/spf13/viper"
)

func TestInitFlow(t *testing.T) {
	t.Run("should be return AuthFlow when bindLocalServer not error", func(t *testing.T) {
		viper.Set(config.AuthHostAddr, "127.0.0.1:8081")
		flow, err := InitFlow()
		flow.server.Close()
		if err != nil {
			t.Errorf("InitFlow() error = %v, wantErr %v", err, false)
			return
		}
	})
	t.Run("should be throw when bindLocalServer error", func(t *testing.T) {
		viper.Set(config.AuthHostAddr, nil)
		_, err := InitFlow()
		if err == nil {
			t.Errorf("InitFlow() error = %v, wantErr %v", err, true)
			return
		}
	})
}

func TestFlow_GetLineWebLoginURL(t *testing.T) {
	t.Run("should be return currect url when set config correctly", func(t *testing.T) {
		viper.Set(config.AuthRedirectURI, "https://example.com/auth")
		viper.Set(config.AuthHost, "https://access.line.me/oauth2/v2.1")
		viper.Set(config.AuthClientID, "1234567890")
		viper.Set(config.AuthScopes, []string{"openid"})
		flow := &Flow{
			state:  "12345abcde",
			nonce:  "09876xyz",
			server: &localServer{},
		}
		want := "https://access.line.me/oauth2/v2.1/authorize?client_id=1234567890&nonce=09876xyz&redirect_uri=https%3A%2F%2Fexample.com%2Fauth&response_type=code&scope=openid&state=12345abcde"
		got, err := flow.GetLineWebLoginURL()
		if err != nil {
			t.Errorf("InitFlow() error = %v, wantErr %v", err, false)
			return
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("RunHealthCheck() = %v, want %v", got, want)
		}

	})
	t.Run("should be throw  when set config incorrectly", func(t *testing.T) {
		viper.Set(config.AuthRedirectURI, nil)
		viper.Set(config.AuthHost, "https://access.line.me/oauth2/v2.1")
		viper.Set(config.AuthClientID, "1234567890")
		viper.Set(config.AuthScopes, []string{"openid"})
		flow := &Flow{
			state:  "12345abcde",
			nonce:  "09876xyz",
			server: &localServer{},
		}
		_, err := flow.GetLineWebLoginURL()
		if err == nil {
			t.Errorf("InitFlow() error = %v, wantErr %v", err, true)
			return
		}

	})
}
