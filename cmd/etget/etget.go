package main

import (
	"flag"
	"log"
	"os"

	"github.com/joneskoo/etget/energiatili"
	"github.com/joneskoo/etget/keyring"
)

func main() {
	credfile := flag.String("credfile", "./credentials.json", "File username/password are saved in (plaintext)")
	output := flag.String("output", "./power.json", "File consumption data is written in")
	flag.Parse()

	cs := keyring.CredentialStore{
		File:   *credfile,
		Domain: "www.energiatili.fi",
	}

	datasource := &energiatili.Client{
		LoginURL:             "https://www.energiatili.fi/Extranet/Extranet/LogIn",
		ConsumptionReportURL: "https://www.energiatili.fi/Reporting/CustomerConsumption/UserConsumptionReport",
		GetUsernamePassword:  cs.UsernamePassword,
	}

	log.Println("Downloading consumption data…")
	var err error
	var fp *os.File
	if fp, err = os.Create(*output); err != nil {
		log.Fatalln(err)
	}
	if err := datasource.ConsumptionReport(fp); err != nil {
		log.Fatalln(err)
	}

	log.Println("OK! Wrote output to", *output)
}
