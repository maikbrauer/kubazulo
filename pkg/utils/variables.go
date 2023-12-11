package utils

import (
	"log"
	"os"
)

const msbaseURL = "https://login.microsoftonline.com/"

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

var (
	CfgClientId         string
	CfgTenantId         string
	CfgForceLogin       string
	CfgLoopbackport     string
	CfgIntermediate     string
	CfgApitokenendpoint string
)

var (
	WarningLogger *log.Logger
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
	DebugLogger   *log.Logger
)

var (
	AuthorizationURL string
	TokenURL         string
)

type Session struct {
	TokenStartTimestamp int64  `json:"tokenstartTimestamp"`
	ExpirationTimestamp int64  `json:"expirationTimestamp"`
	AccessToken         string `json:"access_token"`
	RefreshToken        string `json:"refresh_token"`
}

func FillVariables() {
	AuthorizationURL = msbaseURL + CfgTenantId + "/oauth2/v2.0/authorize"
	TokenURL = msbaseURL + CfgTenantId + "/oauth2/v2.0/token"
}

func init() {
	var logpath = GetHomeDir() + "/.kube/kubazulo/"
	CreateDirectory(logpath)

	file, err := os.OpenFile(logpath+"application.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	InfoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarningLogger = log.New(file, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	DebugLogger = log.New(file, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
}
