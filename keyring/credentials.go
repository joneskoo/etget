// Package keyring implements a poor man's credential store backed by JSON file.
// Using it will store passwords plaintext on disk. Probably don't use it.
package keyring

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"golang.org/x/term"
)

// CredentialStore stores secrets insecurely in a file
type CredentialStore struct {
	// File is the name of the file where credentials are stored
	File string

	// Domain is used in the prompt to user when asking for username
	Domain string

	// return cached value only once
	cacheReturned bool
}

// UsernamePassword returns username and password from cache or prompts user.
// After first call, it will assume cache is invalid and prompt for new
// credentials, storing new value in cache.
func (c *CredentialStore) UsernamePassword() (username, password string, err error) {
	if !c.cacheReturned {
		c.cacheReturned = true
		username, password, err = readCachedCredentials(c.File)
		if err == nil {
			return username, password, nil
		}
	}

	log.Println("Please enter credentials, I will remember them")
	if username, password, err = promptCredentials(c.Domain); err != nil {
		return "", "", err
	}
	log.Println("Storing credentials to", c.File)
	if err = writeCachedCredentials(c.File, username, password); err != nil {
		return "", "", err
	}
	return username, password, nil
}

const jsonIndent = "    "

type credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func promptCredentials(domain string) (username, password string, err error) {
	fmt.Printf("%s username: ", domain)
	if _, err = fmt.Scanln(&username); err != nil {
		return
	}

	fmt.Println(username, " password:")
	passwordBytes, err := term.ReadPassword(0)
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

	if b, err = os.ReadFile(filename); err != nil {
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
	if err = os.WriteFile(filename, b, 0600); err != nil {
		return
	}
	return
}
