package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"golang.org/x/crypto/ssh/terminal"
)

type credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func promptCredentials() (username, password string, err error) {
	fmt.Print("www.energiatili.fi username: ")
	if _, err = fmt.Scanln(&username); err != nil {
		return
	}

	fmt.Println(username, " password:")
	passwordBytes, err := terminal.ReadPassword(0)
	if err != nil {
		return
	}
	password = string(passwordBytes)
	return
}

func readCachedCredentials(filename string) (username, password string, err error) {
	var (
		c credentials
		b []byte
	)

	if b, err = ioutil.ReadFile(filename); err != nil {
		return
	}

	if err = json.Unmarshal(b, &c); err != nil {
		return
	}

	username = c.Username
	password = c.Password
	return
}

func writeCachedCredentials(filename, username, password string) (err error) {
	var b []byte
	c := credentials{username, password}
	if b, err = json.MarshalIndent(c, "", jsonIndent); err != nil {
		return
	}
	if err = ioutil.WriteFile(filename, b, 0600); err != nil {
		return
	}
	return
}
