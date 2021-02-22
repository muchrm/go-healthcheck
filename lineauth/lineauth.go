package lineauth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/muchrm/go-healthcheck/config"
	"github.com/muchrm/go-healthcheck/utils"
	"github.com/pkg/browser"
	"github.com/spf13/viper"
)

// DetectAndAskForToken detect client secret and accesstoken and ask to request a new token
func DetectAndAskForToken() {
	clientSecret, errClientSecret := config.GetString(config.AuthClientSecret)
	clientID, errClientID := config.GetString(config.AuthClientID)
	lineAPIToken, errLineAPIToken := config.GetString(config.AccessToken)
	canStartCalbackServer := errClientSecret == nil && clientSecret != "" && errClientID == nil && clientID != ""

	if (errLineAPIToken != nil || lineAPIToken == "") &&
		!canStartCalbackServer {
		confirm := false
		err := survey.AskOne(&survey.Confirm{
			Message: "Line API token missing and cannot start authen server are you sure to continue?",
		}, &confirm)

		if err != nil {
			fmt.Printf("could not prompt: %s\n sending report to server will ignore", err)
			return
		}
		if confirm {
			return
		}
		os.Exit(0)
	} else if errLineAPIToken != nil || lineAPIToken == "" {
		confirm := false
		err := survey.AskOne(&survey.Confirm{
			Message: "Line API token missing are you want to ask for new token?",
		}, &confirm)
		if err != nil {
			fmt.Printf("could not prompt: %s\n sending report to server will ignore", err)
			return
		}
		if !confirm {
			return
		}
		flow, err := InitFlow()
		if err != nil {
			fmt.Printf("DetectAndAskForToken InitFlow error: %s\nsending report to server will ignore", err)
			return
		}
		token, err := flow.DetectFlow()
		if err != nil {
			fmt.Printf("DetectAndAskForToken DetectFlow error: %s\nsending report to server will ignore", err)
			return
		}

		viper.Set(config.AccessToken, token)
	}

}

type BrowserParams struct {
	ClientID    string
	RedirectURI string
	Scopes      []string
}

type Flow struct {
	server   *localServer
	clientID string
	state    string
	nonce    string
}
type RequestTokenResponse struct {
	AccessToken string `json:"access_token"`
}

// GetLineWebLoginURL generate line url to request authen
func (flow *Flow) GetLineWebLoginURL() (string, error) {
	redirectURI, err := config.GetString(config.AuthRedirectURI)
	if err != nil {
		return "", fmt.Errorf("DetectFlow AuthRedirectURI error:  %w", err)
	}

	ru, err := url.Parse(redirectURI)
	if err != nil {
		return "", fmt.Errorf("DetectFlow redirectURI error:  %w", err)
	}

	host, err := config.GetString(config.AuthHost)
	if err != nil {
		return "", fmt.Errorf("DetectFlow AuthHost error:  %w", err)
	}

	clientID, err := config.GetString(config.AuthClientID)
	if err != nil {
		return "", fmt.Errorf("DetectFlow AuthClientID error:  %w", err)
	}
	scopes, err := config.GetStringArray(config.AuthScopes)
	if err != nil {
		return "", fmt.Errorf("DetectFlow AuthScopes error:  %w", err)
	}

	baseURL := fmt.Sprintf("%s/authorize", host)
	flow.server.CallbackPath = ru.Path
	flow.clientID = clientID

	q := url.Values{}
	q.Set("response_type", "code")
	q.Set("client_id", clientID)
	q.Set("redirect_uri", ru.String())
	q.Set("scope", strings.Join(scopes, " "))
	q.Set("state", flow.state)
	q.Set("nonce", flow.nonce)

	return fmt.Sprintf("%s?%s", baseURL, q.Encode()), nil
}

// DetectFlow start local server to receive code from line authen
func (flow *Flow) DetectFlow() (string, error) {
	go func() {
		_ = flow.server.Serve()
	}()
	defer flow.server.Close()

	browserURL, err := flow.GetLineWebLoginURL()
	if err != nil {
		return "", fmt.Errorf("DetectFlow  GetLineWebLoginURL error:  %w", err)
	}

	err = browser.OpenURL(browserURL)
	if err != nil {
		fmt.Printf("Failed opening a web browser at %s\n,   Please try entering the URL in your browser manually\n", browserURL)
	}

	code := flow.server.WaitForCode()
	if err != nil {
		return "", fmt.Errorf("DetectFlow  WaitForCode error:  %w", err)
	}

	token, err := flow.GetAccessToken(code.Code)
	if err != nil {
		return "", fmt.Errorf("DetectFlow  GetAccessToken error:  %w", err)
	}
	return token, nil
}

// GetAccessToken receive code and use code to ask access token from line authen
func (flow *Flow) GetAccessToken(code string) (string, error) {
	redirectURI, err := config.GetString(config.AuthRedirectURI)
	if err != nil {
		return "", fmt.Errorf("GetAccessToken AuthRedirectURI error:  %w", err)
	}

	clientID, err := config.GetString(config.AuthClientID)
	if err != nil {
		return "", fmt.Errorf("GetAccessToken AuthClientID error:  %w", err)
	}
	clientSecret, err := config.GetString(config.AuthClientSecret)
	if err != nil {
		return "", fmt.Errorf("GetAccessToken AuthClientSecret error:  %w", err)
	}

	host, err := config.GetString(config.AuthResourceHose)
	if err != nil {
		return "", fmt.Errorf("GetAccessToken AuthHost error:  %w", err)
	}

	baseURL := fmt.Sprintf("%s/token", host)

	body := url.Values{}
	body.Set("grant_type", "authorization_code")
	body.Set("code", code)
	body.Set("redirect_uri", redirectURI)
	body.Set("client_id", clientID)
	body.Set("client_secret", clientSecret)

	req, err := http.NewRequest(
		http.MethodPost,
		baseURL,
		strings.NewReader(body.Encode()),
	)
	if err != nil {
		return "", fmt.Errorf("GetAccessToken NewRequest error:  %w", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(body.Encode())))
	client := http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("GetAccessToken Do error:  %w", err)
	}
	defer resp.Body.Close()

	var respBody RequestTokenResponse
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		return "", fmt.Errorf("GetAccessToken Decode error:  %w", err)
	}

	accessToken := respBody.AccessToken
	if accessToken == "" {
		return "", fmt.Errorf("GetAccessToken AccessToken not found")
	}
	return accessToken, nil
}

// InitFlow create a local server instance and oauth flow
func InitFlow() (*Flow, error) {
	server, err := bindLocalServer()
	if err != nil {
		return nil, fmt.Errorf("InitFlow error %w", err)
	}

	state, _ := utils.RandomString(20)
	nonce, _ := utils.RandomString(20)

	return &Flow{
		server: server,
		state:  state,
		nonce:  nonce,
	}, nil
}
