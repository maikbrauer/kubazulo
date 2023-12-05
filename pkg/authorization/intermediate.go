package authorization

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"kubazulo/pkg/utils"
	"log"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

type Jsondata struct {
	AuthCode    string `json:"authcode"`
	RedirectURI string `json:"redirect_uri"`
}

func GetTokenData(data Jsondata) (t Tokens, err error) {
	marshalled, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("impossible to marshall teacher: %s", err)
	}

	req, err := http.NewRequest("POST", utils.Cfg_apitokenendpoint, bytes.NewReader(marshalled))
	req.Header.Set("Content-Type", "application/json")
	client := http.Client{Timeout: 10 * time.Second}
	response, err := client.Do(req)

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
