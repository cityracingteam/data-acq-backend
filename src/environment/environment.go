package environment

import (
	"fmt"
	"os"
)

func GetEnvOrDefault(key string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	fmt.Println("[warning]: Using default value for env key " + key)
	return devEnv[key]
}

func GetCallbackUri() string {
	return (GetEnvOrDefault("ENDPOINT_SCHEME") +
		"://" +
		GetEnvOrDefault("ENDPOINT") +
		"/auth/openid-connect/callback")
}

var devEnv = map[string]string{
	"DOMAIN":                       "localhost",
	"ENDPOINT":                     "localhost:8080",
	"ENDPOINT_SCHEME":              "http",
	"OPENID_CONNECT_DISCOVERY_URL": "http://localhost:8081/.well-known/openid-configuration",
	"OPENID_CONNECT_KEY":           "backend",
	"OPENID_CONNECT_SECRET":        "development-secret",
}
