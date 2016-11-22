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

const jsonIndent = "    "

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

type ConsumptionFetcher interface {
	Login(username, password string) error
	ConsumptionReport(w io.Writer) error
}

func main() {
	credfile := flag.String("credfile", "./credentials.json", "File username/password are saved in (plaintext)")
	output := flag.String("output", "./power.json", "File consumption data is written in")
	flag.Parse()

	username, password, err := readCachedCredentials(*credfile)
	if err != nil {
		log.Println("Could not find cached credentials in", *credfile)
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

	log.Println("Storing credentials to", *credfile)
	if err = writeCachedCredentials(*credfile, username, password); err != nil {
		log.Fatalln(err)
	}

	log.Println("Downloading consumption data…")
	var fp *os.File
	if fp, err = os.Create(*output); err != nil {
		log.Fatalln(err)
	}
	if err := f.ConsumptionReport(fp); err != nil {
		log.Fatalln(err)
	}

	log.Println("OK! Wrote output to", *output)
}
