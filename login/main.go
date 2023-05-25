package main

import (
	"log"
	"os"
)

func getEnv(key string) (env string) {
	env = os.Getenv(key)
	if env == "" {
		log.Fatalf("Env %s not set", key)
		return
	}
	return
}

func main() {
	ntu := NoTypeUsername{
		CredPath: getEnv("CRED_PATH"),
		SamlPath: getEnv("SAML_PATH"),
	}
	ntu.LoginNTUCOOL()
}
