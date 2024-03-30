package authorization

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"kubazulo/pkg/utils"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

func AuthorizeRequestDeviceFlow(c utils.AuthorizationConfig) (jsonDf jsonDeviceFlow) {

	formVals := url.Values{}
	formVals.Set("client_id", c.ClientID)
	formVals.Set("scope", utils.CfgClientId+"/.default offline_access")

	uri, _ := url.Parse(utils.AuthorizationURLDevice)
	uri.RawQuery = formVals.Encode()

	req, err := http.NewRequest(http.MethodGet, uri.String(), nil)
	req.Header.Set("Content-Type", "application/json")
	client := http.Client{Timeout: 10 * time.Second}
	response, err := client.Do(req)

	if err != nil {
		errors.Wrap(err, "error while trying to get tokens")
	}
	body, err := io.ReadAll(response.Body)

	if err != nil {
		errors.Wrap(err, "error while trying to read token json body")
	}

	err = json.Unmarshal(body, &jsonDf)
	if err != nil {
		errors.Wrap(err, "error while trying to parse token json body")
	}

	fmt.Fprintln(os.Stderr, "Please visit", jsonDf.VerificationUri, "and put the following code there:", jsonDf.UserCode)
	fmt.Fprintln(os.Stderr, "You have 3 Minutes to complete !!")

	if err != nil {
		panic(errors.Wrap(err, "Error while opening DeviceFlow URL"))
	}

	return
}

func AuthorizeRequest(c utils.AuthorizationConfig) (token AuthorizationCode) {
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

	err := open(uri.String())
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

func open(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
		url = strings.ReplaceAll(url, "&", "^&")
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}

	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}
