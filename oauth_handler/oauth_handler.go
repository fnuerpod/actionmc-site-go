package oauth_handler

import (
	"errors"

	"io"

	"git.dsrt-int.net/actionmc/actionmc-site-go/config"
	"git.dsrt-int.net/actionmc/actionmc-site-go/logging"
	"github.com/buger/jsonparser"
	"github.com/pollen5/discord-oauth2"
	"golang.org/x/oauth2"
)

// main OAuth handler struct.
// will use struct methods instead of relying on global variables.

type OAuthHandler struct {
	oauth_config       *oauth2.Config
	logger             *logging.Logger
	site_configuration *config.Configuration
}

// defines a user instead of having to define them as individual
// variables. much more cleaner and makes code more legible.
type OAuth_User struct {
	Username  string
	DiscordId string
}

func ParseUserJSON(body []byte) (error, *OAuth_User) {
	uname, uname_err := jsonparser.GetString(body, "username")
	id, id_err := jsonparser.GetString(body, "id")

	if (id_err != nil) || (uname_err != nil) {
		return errors.New("Failed to parse user JSON."), new(OAuth_User)
	}

	return nil, &OAuth_User{
		Username:  uname,
		DiscordId: id,
	}

}

func GetConfig(site_configuration *config.Configuration) *oauth2.Config {
	if site_configuration.IsPreproduction {
		return &oauth2.Config{
			Endpoint:     discord.Endpoint,
			Scopes:       []string{discord.ScopeIdentify},
			RedirectURL:  "https://actionmc.ml/login/callback",
			ClientID:     "800316621343162429",
			ClientSecret: "JCYm4OAHd7Ygm_zBZHdvgb3dbv1GMMm0",
		}
	} else if site_configuration.IsBeta {
		return &oauth2.Config{
			Endpoint:     discord.Endpoint,
			Scopes:       []string{discord.ScopeIdentify},
			RedirectURL:  "https://b.actionmc.ml/login/callback",
			ClientID:     "800316621343162429",
			ClientSecret: "JCYm4OAHd7Ygm_zBZHdvgb3dbv1GMMm0",
		}
	} else {
		return &oauth2.Config{
			Endpoint:     discord.Endpoint,
			Scopes:       []string{discord.ScopeIdentify},
			RedirectURL:  "https://actionmc.ml/login/callback",
			ClientID:     "800316621343162429",
			ClientSecret: "JCYm4OAHd7Ygm_zBZHdvgb3dbv1GMMm0",
		}
	}
}

func New() *OAuthHandler {
	// create a new oauth handler
	conf := config.InitialiseConfig()

	return &OAuthHandler{
		oauth_config:       GetConfig(conf),
		logger:             logging.New(),
		site_configuration: conf,
	}
}

// struct methods for OAuth Handler

func (OH *OAuthHandler) GetConf() *oauth2.Config {
	return OH.oauth_config
}

func (OH *OAuthHandler) GetInformation(access_token string) (error, *OAuth_User) {
	// attempt to get the discord query token using the given access token.
	token, err := OH.oauth_config.Exchange(oauth2.NoContext, access_token)

	//OH.logger.Fatal.Fatalln(err)

	if err != nil {
		// failed to get query token via access token.
		return errors.New("Failed to get query token from access token."), new(OAuth_User)
	}

	// try using our given query token to get the user information from
	// discord.

	res, err := OH.oauth_config.Client(oauth2.NoContext, token).Get("https://discordapp.com/api/users/@me")

	if err != nil || res.StatusCode != 200 {
		// query token failed
		return errors.New("Failed to get user information using query token."), new(OAuth_User)
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if err != nil {
		// reading from body failed. this is kinda major.
		return errors.New("Failed to read user information from body."), new(OAuth_User)
	}

	err, user := ParseUserJSON(body)

	if err != nil {
		return err, user
	}

	return nil, user
}
