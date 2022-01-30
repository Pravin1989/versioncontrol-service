package services

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"versioncontrol-service/src/config"
	"versioncontrol-service/src/models"

	"github.com/google/go-github/v42/github"
	"golang.org/x/oauth2"
)

var (
	oauthConfig = &oauth2.Config{
		ClientID:     "",
		ClientSecret: "",
		RedirectURL:  "",
		Scopes:       []string{},
		Endpoint:     oauth2.Endpoint{},
	}
)

// This method configures the OAuth Server
func InitializeOAuthGithub() {
	oauthConfig.ClientID = config.OAuthConfig.ClientID
	oauthConfig.ClientSecret = config.OAuthConfig.ClientSecret
	oauthConfig.Endpoint = oauth2.Endpoint{AuthURL: config.OAuthConfig.AuthURL, TokenURL: config.OAuthConfig.TokenURL}
	oauthConfig.Scopes = config.OAuthConfig.Scopes
	oauthConfig.RedirectURL = config.ServerConfig.HostName + ":" + strconv.Itoa(config.ServerConfig.Port) + config.OAuthConfig.RedirectURL
}

// This method is works as starting point of the service
func HandleHomePage(w http.ResponseWriter, r *http.Request) {
	handleLogin(w, r, oauthConfig, config.OAuthConfig.OauthStateString)
}
func handleLogin(w http.ResponseWriter, r *http.Request, oauthConf *oauth2.Config, oauthStateString string) {
	URL, err := url.Parse(oauthConf.Endpoint.AuthURL)
	if err != nil {
		log.Println("Parse failed : " + err.Error())
		return
	}
	log.Println("Authorized URL : ", URL.String())
	parameters := url.Values{}
	parameters.Add("client_id", oauthConf.ClientID)
	parameters.Add("scope", strings.Join(oauthConf.Scopes, " "))
	parameters.Add("redirect_uri", oauthConf.RedirectURL)
	parameters.Add("response_type", "code")
	parameters.Add("state", oauthStateString)
	URL.RawQuery = parameters.Encode()
	url := URL.String()
	log.Println("Redirect to Github for Access Code : ", url)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

//THis method is call back handler for OAuth Autharization
func HandleCallBackFromGithubAuth(w http.ResponseWriter, r *http.Request) {
	log.Println("Callback received from Github")
	code := r.FormValue("code")
	if code == "" {
		log.Println("Code not found..")
		w.Write([]byte("Code Not Found to provide AccessToken..\n"))
		reason := r.FormValue("error_reason")
		if reason == "user_denied" {
			w.Write([]byte("User has denied Permission.."))
		}
	} else {
		token, err := oauthConfig.Exchange(oauth2.NoContext, code)
		if err != nil {
			log.Println("Exchange failed with error : " + err.Error() + "\n")
			return
		}
		ctx := context.Background()
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token.AccessToken},
		)
		tc := oauth2.NewClient(ctx, ts)

		client := github.NewClient(tc)

		models.StartProcess(ctx, client)
	}
}
