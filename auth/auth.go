package auth

import (
	"net/http"
	"strings"

	"github.com/rbrick/imgy/config"
	"github.com/rbrick/imgy/db"
	"github.com/rbrick/imgy/util"
	"golang.org/x/oauth2"
)

var registeredServices = []Service{}
var enabledServices = []Service{}

var oauthUrlPath string

// UserInfo is the user information we received from authentication services
type UserInfo struct {
	Email, DisplayName, Picture string
}

type BaseAuthService struct {
	config *oauth2.Config
	path   string
}

func (bas *BaseAuthService) Path() string {
	return bas.path
}

func (bas *BaseAuthService) Config() *oauth2.Config {
	return bas.config
}

func (bas *BaseAuthService) AuthURL(state string) string {
	return bas.config.AuthCodeURL(state)
}

// Service represents a service that provides authentication.
type Service interface {
	Name() string
	Path() string
	Setup(*config.OAuthConfig) *oauth2.Config
	Config() *oauth2.Config
	AuthURL(string) string
	Callback(*http.Client, *oauth2.Token) (*UserInfo, error)
}

// RegisterService registers a service
func RegisterService(service Service) {
	registeredServices = append(registeredServices, service)
}

// Services get all the services that are enabled
func Services() []Service {
	return enabledServices
}

// GetService gets a service by it's name
func GetService(name string) Service {
	for _, v := range registeredServices {
		if strings.EqualFold(v.Name(), name) {
			return v
		}
	}
	return nil
}

func Init(c *config.Config) {
	oauthUrlPath = c.OauthURL
	for _, name := range c.OauthProviders {
		service := GetService(name)
		service.Setup(c.OauthConfigs[name])
		enabledServices = append(enabledServices, service)
	}
}

func OAuthCallbackHandler(service Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		state := r.URL.Query().Get("state")

		sess := util.MustSession(r, "imgy-auth")

		if state != sess.Values["state"].(string) {
			http.Error(w, "Invalid state", http.StatusUnauthorized)
		} else {
			token, err := service.Config().Exchange(oauth2.NoContext, r.URL.Query().Get("code"))
			if err != nil {
				panic(err)
			}

			if token.Valid() {
				client := service.Config().Client(oauth2.NoContext, token)

				// Get the users email
				userInfo, err := service.Callback(client, token)

				if err != nil {
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}

				imgySess := util.MustSession(r, "imgy")

				if u := db.GetUserByEmail(userInfo.Email); u != nil {
					// Update the profile picture in case it changed
					u.ProfilePicture = userInfo.Picture
					u.StartSession(imgySess, r, w)
				} else {
					u = &db.User{
						UserID:         util.GetRandom(8),
						DisplayName:    userInfo.DisplayName,
						Email:          userInfo.Email,
						ProfilePicture: userInfo.Picture,
						UploadToken:    util.GetRandom(32),
					}
					u.StartSession(imgySess, r, w)
				}

				http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			}
		}
	}
}
