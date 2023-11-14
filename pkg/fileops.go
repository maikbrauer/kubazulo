package kubazulo

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/pkg/errors"
)

func WriteSession(Expiry int64, TokenStart int64, _AccessToken string, _RefreshToken string) {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	path := home + "/.kube/cache/kubazulo"
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			log.Println(err)
		}
	}

	f, err := os.Create(home + "/.kube/cache/kubazulo/azuredata.json")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	data := Session{
		TokenStartTimestamp: TokenStart,
		ExpirationTimestamp: Expiry,
		AccessToken:         _AccessToken,
		RefreshToken:        _RefreshToken,
	}

	a_json, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	f.Write(a_json)
}

func ReadSession() Session {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	fileContent, err := os.Open(home + "/.kube/cache/kubazulo/azuredata.json")

	if err != nil {
		log.Fatal(err)
	}

	defer fileContent.Close()

	byteResult, _ := ioutil.ReadAll(fileContent)

	data := Session{}

	err = json.Unmarshal([]byte(byteResult), &data)

	return data
}
