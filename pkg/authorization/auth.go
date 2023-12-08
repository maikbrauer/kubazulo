package authorization

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"kubazulo/pkg/utils"
	"log"
	"net/http"
	"net/url"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
)

const SuccessMsg = `
<!DOCTYPE html> <html lang="en"> 
<head>
<style>
body {
  background-image: url('https://wallpaperaccess.com/full/4834955.jpg');
  background-repeat: no-repeat;
  background-attachment: fixed;
  background-size: 100% 100%;
}
</style>
<meta charset="UTF-8">
    <title>Azure Platform Authentication Service</title>
</head>
<body>
	<p style="background-image: url('https://wallpaperaccess.com/full/4834955.jpg');"></p>
    <h2><p style="color: white">You have been successfully authenticated and now ready to communicate with the API-Server</br></br>
	You can close the Browser window now and get back to the command-line!</p></h2>
</body>
</html>
`

// AuthorizationCode is a value provided after initial successful
// authentication/authorization, it is used to get access/refresh tokens
type AuthorizationCode struct {
	Value string
}

// Tokens holds access and refresh tokens
type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Expiry       int    `json:"expires_in"`
}

// GetTokens retrieves access and refresh tokens for a given scope
func GetTokens(c utils.AuthorizationConfig, authCode AuthorizationCode, scope string) (t Tokens, err error) {
	formVals := url.Values{}
	formVals.Set("code", authCode.Value)
	formVals.Set("grant_type", "authorization_code")
	formVals.Set("redirect_uri", c.RedirectURL())
	formVals.Set("scope", scope)
	if c.ClientSecret != "" {
		formVals.Set("client_secret", c.ClientSecret)
	}
	formVals.Set("client_id", c.ClientID)
	response, err := http.PostForm(utils.TokenURL, formVals)

	if err != nil {
		return t, errors.Wrap(err, "error while trying to get tokens")
	}
	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return t, errors.Wrap(err, "error while trying to read token json body")
	}

	err = json.Unmarshal(body, &t)
	if err != nil {
		return t, errors.Wrap(err, "error while trying to parse token json body")
	}

	return
}

// startLocalListener opens a http server to retrieve the redirect from initial
// authentication and set the authorization code's value
func startLocalListener(c utils.AuthorizationConfig, token *AuthorizationCode) *http.Server {
	srv := &http.Server{Addr: fmt.Sprintf(":%s", c.RedirectPort)}

	http.HandleFunc(c.RedirectPath, func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			log.Fatalf("Error while parsing form from response %s", err)
			return
		}
		for k, v := range r.Form {
			if k == "code" {
				token.Value = strings.Join(v, "")
			}
		}

		fmt.Fprintf(w, "%s", SuccessMsg)
	})

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			// cannot panic, because this probably is an intentional close
			//log.Printf("Httpserver: ListenAndServe() error: %s", err)
		}
	}()

	// returning reference so caller can call Shutdown()
	return srv
}

// LoginRequest asks the os to open the login URL and starts a listening on the
// configured port for the authorization code. This is used on initial login to
// get the initial token pairs
func LoginRequest(c utils.AuthorizationConfig) (token AuthorizationCode) {
	formVals := url.Values{}
	formVals.Set("grant_type", "authorization_code")
	formVals.Set("redirect_uri", c.RedirectURL())
	formVals.Set("scope", utils.CfgClientId+"/.default")
	formVals.Set("response_type", "code")
	formVals.Set("response_mode", "query")
	formVals.Set("client_id", c.ClientID)
	formVals.Set("state", "12345")
	if strings.ToLower(utils.CfgForceLogin) == "true" {
		formVals.Set("prompt", "login")
	}
	uri, _ := url.Parse(utils.AuthorizationURL)
	uri.RawQuery = formVals.Encode()

	cmd := exec.Command(c.OpenCMD, uri.String())
	err := cmd.Start()
	if err != nil {
		panic(errors.Wrap(err, "Error while opening login URL"))

	}
	running := true
	srv := startLocalListener(c, &token)
	for running {
		if token.Value != "" {
			if err := srv.Shutdown(context.TODO()); err != nil {
				panic(err) // failure/timeout shutting down the server gracefully
			}
			running = false
		}
	}
	return
}

func RenewAccessToken(refreshToken string) (t Tokens, err error) {
	if utils.CfgIntermediate == "true" {

		var data = JsonData{}
		data.GrantType = "refresh_token"
		data.Code = refreshToken

		t, err := GetTokenData(data)
		return t, err

	} else {

		formVals := url.Values{}
		formVals.Set("refresh_token", refreshToken)
		formVals.Set("grant_type", "refresh_token")
		formVals.Set("client_id", utils.CfgClientId)
		response, err := http.PostForm(utils.TokenURL, formVals)

		if err != nil {
			return t, errors.Wrap(err, "error while trying to get tokens")
		}
		body, err := ioutil.ReadAll(response.Body)

		if err != nil {
			return t, errors.Wrap(err, "error while trying to read token json body")
		}

		err = json.Unmarshal(body, &t)
		if err != nil {
			return t, errors.Wrap(err, "error while trying to parse token json body")
		}

		return t, err
	}
}
