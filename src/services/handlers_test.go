package services

import (
	"net/http"
	"testing"
	"versioncontrol-service/src/config"
	"versioncontrol-service/src/tests"

	"github.com/gorilla/mux"
)

var (
	router *mux.Router
)

func Test_HandleHomePage(t *testing.T) {
	router = startTestServer()
	url := "https://localhost:8080/versioncontrol-service/"
	t.Run("Redirect to Github for Authentication", func(t *testing.T) {
		tests.ExecuteAndParseJSON(t, router, http.MethodGet, "/", "", http.StatusTemporaryRedirect, url)
	})
}
func TestInitializeOAuthGithub(t *testing.T) {
	config.OAuthConfig.ClientID = "TestId"
	config.OAuthConfig.ClientSecret = "TestSecret"
	config.OAuthConfig.AuthURL = ""
	config.OAuthConfig.RedirectURL = "http://test.com/versioncontrol-service/callback"
	config.OAuthConfig.Scopes = []string{"repo"}
	tests := []struct {
		name string
	}{
		{"Success : Configure OAuth Details"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			InitializeOAuthGithub()
		})
	}
}

func startTestServer() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", HandleHomePage)
	return r
}
