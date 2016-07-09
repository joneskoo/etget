package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/joneskoo/etget/fetcher"
	"golang.org/x/crypto/ssh/terminal"
)

var (
	// CredentialsFile is the file name where credentials are stored
	CredentialsFile = flag.String("credfile", "./credentials.json",
		"File username/password are saved in (plaintext)")

	// OutputFile is the name of the file where we write output
	OutputFile = flag.String("output", "./power.json", "File consumption data is written in")
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
	if b, err = json.Marshal(c); err != nil {
		return
	}
	if err = ioutil.WriteFile(filename, b, 0600); err != nil {
		return
	}
	return
}

type ConsumptionFetcher interface {
	Login(username, password string) error
	ConsumptionReport(w io.Writer) error
}

func main() {
	flag.Parse()
	username, password, err := readCachedCredentials(*CredentialsFile)
	if err != nil {
		log.Println("Could not find cached credentials in", *CredentialsFile)
		log.Println("Please enter credentials, I will remember them")
		if username, password, err = promptCredentials(); err != nil {
			log.Fatalln(err)
		}
	}

	var f ConsumptionFetcher = &fetcher.Fetcher{}

	log.Println("Logging in…")
	if err = f.Login(username, string(password)); err != nil {
		log.Println("Could not log in:", err)
		username, password, err = promptCredentials()
		if err = f.Login(username, string(password)); err != nil {
			log.Fatalln(err)
		}
	}
	log.Println("Login OK.")

	log.Println("Storing credentials to", CredentialsFile)
	if err = writeCachedCredentials(*CredentialsFile, username, password); err != nil {
		log.Fatalln(err)
	}

	log.Println("Downloading consumption data…")
	var fp *os.File
	if fp, err = os.Create(*OutputFile); err != nil {
		log.Fatalln(err)
	}
	if err := f.ConsumptionReport(fp); err != nil {
		log.Fatalln(err)
	}

	log.Println("OK! Wrote output to", OutputFile)
}
