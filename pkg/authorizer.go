package kubazulo

import (
	"fmt"
	"net/url"
)

type AuthorizationConfig struct {
	Host         string
	Scheme       string
	RedirectPort string
	RedirectPath string
	Scope        string
	ClientID     string
	OpenCMD      string
	ClientSecret string
}

var DefaultConfig = AuthorizationConfig{
	Host:         "localhost",
	Scheme:       "http",
	RedirectPath: "/",
	Scope:        "openid profile offline_access user.read",
	OpenCMD:      "open",
}

// RedirectURL
func (c AuthorizationConfig) RedirectURL() string {
	if c.RedirectPort != "" {
		c.Host = fmt.Sprintf("%s:%s", c.Host, c.RedirectPort)
	}

	uri := url.URL{
		Host:   c.Host,
		Scheme: c.Scheme,
		Path:   c.RedirectPath,
	}

	return uri.String()
}
