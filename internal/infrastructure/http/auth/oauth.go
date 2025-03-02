package auth

import (
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
	"golang.org/x/oauth2"
	"net/http"
	"store-manager/internal/infrastructure/config"
)

type OAuthConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
	Scopes       []string
	Endpoint     oauth2.Endpoint
}

func InitGoth(OAuthConfig OAuthConfig) {
	store := sessions.NewCookieStore([]byte(config.EnvConfigs.SessionSecret))
	store.Options = &sessions.Options{
		Path:     "/",
		HttpOnly: true,
		MaxAge:   config.EnvConfigs.MaxAge,
		Secure:   config.EnvConfigs.Prod,
		SameSite: http.SameSiteLaxMode,
	}

	gothic.Store = store

	goth.UseProviders(
		google.New(OAuthConfig.ClientID, OAuthConfig.ClientSecret, OAuthConfig.RedirectURI, OAuthConfig.Scopes...),
	)
}
