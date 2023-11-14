package kubazulo

import (
	"fmt"
	"net/url"
)

type AuthorizationConfig struct {
	RedirectPort string
	RedirectPath string
	Scope        string
	ClientID     string
	OpenCMD      string
	ClientSecret string
}

var DefaultConfig = AuthorizationConfig{
	RedirectPort: "8080",
	RedirectPath: "/",
	Scope:        "openid profile offline_access user.read",
	OpenCMD:      "open",
}

// RedirectURL
func (c AuthorizationConfig) RedirectURL() string {
	host := "localhost"
	if c.RedirectPort != "" {
		host = fmt.Sprintf("%s:%s", host, c.RedirectPort)
	}
	uri := url.URL{
		Host:   host,
		Scheme: "http",
		Path:   c.RedirectPath,
	}

	return uri.String()
}
