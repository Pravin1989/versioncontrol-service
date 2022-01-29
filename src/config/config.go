package config

import (
	"encoding/json"
	"errors"
	"os"
)

var (
	//ServerConfig stores application configuration
	ServerConfig serverConfiguration
	//OAuthConfig stores OAuth Configurations
	OAuthConfig oAuthConfiguration
)

type oAuthConfiguration struct {
	ClientID         string   `json:"clientId"`
	ClientSecret     string   `json:"clientSecret"`
	OauthStateString string   `json:"oauthStateString"`
	AuthURL          string   `json:"authUrl"`
	TokenURL         string   `json:"tokenUrl"`
	RedirectURL      string   `json:"redirectUrl"`
	Scopes           []string `json:"scopes"`
}

type serverConfiguration struct {
	HostName string `json:"hostName"`
	Port     int    `json:"port"`
}

func loadAppConfig() error {
	bytes := []byte(os.Getenv("APP_SERVER_PROPERTIES"))
	if err := json.Unmarshal(bytes, &ServerConfig); err != nil {
		return err
	}
	if ServerConfig.Port == 0 {
		return errors.New("Mandatory configuration missing for application config")
	}
	return nil
}

func loadOAuthServerConfig() error {
	bytes := []byte(os.Getenv("OAuth_SERVER_PROPERTIES"))
	if err := json.Unmarshal(bytes, &OAuthConfig); err != nil {
		return err
	}
	if OAuthConfig.ClientID == "" || OAuthConfig.ClientSecret == "" || OAuthConfig.AuthURL == "" || OAuthConfig.TokenURL == "" || OAuthConfig.RedirectURL == "" || len(OAuthConfig.Scopes) < 1 {
		return errors.New("Mandatory configuration missing for OAuth config")
	}
	return nil
}
func Load() error {
	resetConfig()
	if err := loadAppConfig(); err != nil {
		return err
	}
	if err := loadOAuthServerConfig(); err != nil {
		return err
	}
	return nil
}
func resetConfig() {
	ServerConfig = serverConfiguration{}
	OAuthConfig = oAuthConfiguration{}
}
