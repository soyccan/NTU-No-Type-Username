package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (c *credentials) valid() bool {
	return len(c.Username) > 0
}

func loadCredentials(path string) (cred credentials, err error) {
	fileData, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal("Error reading JSON file:", err)
		return
	}

	err = json.Unmarshal(fileData, &cred)
	if err != nil || !cred.valid() {
		log.Fatal("Error parsing JSON data:", err)
		return
	}
	return
}
