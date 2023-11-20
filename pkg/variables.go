package kubazulo

import (
	"log"
	"os"
)

var (
	Cfg_client_id    string
	Cfg_tenant_id    string
	Cfg_force_login  string
	Cfg_loopbackport string
)

var (
	WarningLogger *log.Logger
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
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
	// AuthorizationURL is the endpoint used for initial login/auth
	AuthorizationURL = "https://login.microsoftonline.com/" + Cfg_tenant_id + "/oauth2/v2.0/authorize"
	// TokenURL is the endpoint for getting access/refresh tokens
	TokenURL = "https://login.microsoftonline.com/" + Cfg_tenant_id + "/oauth2/v2.0/token"
}

func init() {
	var logpath string = GetHomeDir() + "/.kube/kubazulo/"
	createDirectory(logpath)

	file, err := os.OpenFile(logpath+"application.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	InfoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarningLogger = log.New(file, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}
