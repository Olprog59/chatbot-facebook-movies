package commons

import (
	"log"
	"os"
)

const (
	GRAPHQL_URL  = "https://graph.facebook.com/v19.0"
	URL_CALLBACK = "https://sacred-certain-lemming.ngrok-free.app"
)

var (
	FB_TOKEN      = getEnv("FB_TOKEN")
	WEBHOOK_TOKEN = getEnv("WEBHOOK_TOKEN")
	NOCO_URL      = getEnv("NOCO_URL")
	NOCO_TABLE_ID = getEnv("NOCO_TABLE_ID")
	NOCO_API_KEY  = getEnv("NOCO_API_KEY")
)

func getEnv(key string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		log.Printf("%s not set\n", key)
		return ""
	}
	return val
}
