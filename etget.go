package main

import (
	"flag"
	"io"
	"log"
	"os"

	"github.com/joneskoo/etget/fetcher"
)

const jsonIndent = "    "

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
