package kubazulo

import (
	"log"
	"os"
	"time"
)

func GetCurrentUnixTime() int64 {
	return time.Now().Unix()
}

func GetExpiryUnixTime(ExpiryWindow int64) int64 {
	return time.Now().Unix() + ExpiryWindow
}

func ConvertUnixToRFC3339(timestamp int64) string {
	return time.Unix(timestamp, 0).UTC().Format("2006-01-02T15:04:05Z07:00")
}

func GetHomeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	return home
}
