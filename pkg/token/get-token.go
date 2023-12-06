package token

import (
	"encoding/json"
	"fmt"
	"kubazulo/pkg/authorization"
	"kubazulo/pkg/utils"
	"log"

	"github.com/pkg/errors"
	"github.com/spf13/pflag"

	"os"
	"time"
)

type Spec struct {
	Interactive bool `json:"interactive"`
}

type Status struct {
	ExpirationTimestamp string `json:"expirationTimestamp"`
	Token               string `json:"token"`
}

type Product struct {
	Kind       string `json:"kind"`
	ApiVersion string `json:"apiVersion"`
	Spec       Spec   `json:"spec"`
	Status     Status `json:"status"`
}

type CliFlag struct {
	Required bool
	Usage    string
	Name     string
	Address  *string
}

func kubeoutput(accesstoken string) {
	kcoutput := Product{
		Kind:       "ExecCredential",
		ApiVersion: "client.authentication.k8s.io/v1",
		Spec:       Spec{false},
		Status:     Status{utils.ConvertUnixToRFC3339(utils.GetCurrentUnixTime()), accesstoken},
	}
	bytes, _ := json.Marshal(kcoutput)
	fmt.Println(string(bytes))
}

func createNewToken() {
	authConfig := utils.DefaultConfig
	authConfig.ClientID = utils.CfgClientId
	authConfig.RedirectPort = utils.CfgLoopbackport
	//authConfig.ClientSecret = x.ClientSecret

	authCode := authorization.LoginRequest(authConfig)
	if utils.CfgIntermediate == "true" {

		var data = authorization.JsonData{
			AuthCode:    authCode.Value,
			RedirectURI: "http://localhost:" + utils.CfgLoopbackport,
		}

		t, err := authorization.GetTokenData(data)
		if err != nil {
			panic(err)
		}
		kubeoutput(t.AccessToken)
		utils.WriteSession(utils.GetExpiryUnixTime(int64(t.Expiry)), utils.GetCurrentUnixTime(), t.AccessToken, t.RefreshToken)
	} else {
		t, err := authorization.GetTokens(authConfig, authCode, "profile openid offline_access")
		if err != nil {
			panic(err)
		}
		kubeoutput(t.AccessToken)
		utils.WriteSession(utils.GetExpiryUnixTime(int64(t.Expiry)), utils.GetCurrentUnixTime(), t.AccessToken, t.RefreshToken)
	}
}

func GetTokenProcess(flags *pflag.FlagSet) {
	var _r utils.Session

	utils.InfoLogger.Println("Application invoked")

	if utils.CheckFlagExistence(flags, "client-id") {
		utils.CfgClientId = flags.Lookup("client-id").Value.String()
	}

	if utils.CheckFlagExistence(flags, "tenant-id") {
		utils.CfgTenantId = flags.Lookup("tenant-id").Value.String()
	}

	if utils.CheckFlagExistence(flags, "force-login") {
		utils.CfgForceLogin = flags.Lookup("force-login").Value.String()
	}

	if utils.CheckFlagExistence(flags, "loopbackport") {
		utils.CfgLoopbackport = flags.Lookup("loopbackport").Value.String()
	}

	if utils.CheckFlagExistence(flags, "intermediate") {
		utils.CfgIntermediate = flags.Lookup("intermediate").Value.String()
	}

	if utils.CheckFlagExistence(flags, "api-token-endpoint") {
		utils.CfgApitokenendpoint = flags.Lookup("api-token-endpoint").Value.String()
	}

	utils.FillVariables()

	if _, err := os.Stat(utils.GetHomeDir() + "/.kube/cache/kubazulo/azuredata.json"); errors.Is(err, os.ErrNotExist) {
		utils.InfoLogger.Println("Cache File does not exist. New AccessToken obtained from Azure-API")
		createNewToken()
	} else {
		r := utils.ReadSession()
		_r = r
		if _r.AccessToken == "" {
			utils.InfoLogger.Println("Cache File does not contain an Access-Token. New AccessToken obtained from Azure-API")
			createNewToken()
		} else if time.Now().Unix() >= _r.ExpirationTimestamp {
			utils.InfoLogger.Println("Cache File exist but AccessToken is expired. New AccessToken obtained from Azure-API via Refreshtoken")
			t, err := authorization.RenewAccessToken(_r.RefreshToken)
			if err != nil {
				log.Fatal(err)
			}
			kubeoutput(t.AccessToken)
			utils.WriteSession(utils.GetExpiryUnixTime(int64(t.Expiry)), utils.GetCurrentUnixTime(), t.AccessToken, t.RefreshToken)
			utils.InfoLogger.Println("Cache File updated with the latest information from Azure-API")
		} else {
			utils.InfoLogger.Println("Cache File exist. AccessToken taken from cache file")
			kubeoutput(r.AccessToken)
		}
	}
}
