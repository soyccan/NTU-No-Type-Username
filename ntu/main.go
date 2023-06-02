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
	ntu := NtuCOOL{
		NoTypeUsername{
			CredPath:   getEnv("CRED_PATH"),
			CookiePath: getEnv("COOKIE_PATH"),
		},
	}
	Login(&ntu)
}
