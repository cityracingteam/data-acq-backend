package environment

import "os"

func GetEnvOrDefault(key string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return devEnv[key]
}

var devEnv = map[string]string{
	"ENDPOINT":                     "http://localhost:8080",
	"OPENID_CONNECT_DISCOVERY_URL": "http://localhost:8081/.well-known/openid-configuration",
	"OPENID_CONNECT_KEY":           "backend",
	"OPENID_CONNECT_SECRET":        "development-secret",
}
