package main

import (
	// "bufio"
	"fmt"
	"os"

	"github.com/joneskoo/etget/fetcher"
	"github.com/joneskoo/etget/secrets"
	// "golang.org/x/crypto/ssh/terminal"
)

func main() {
	// reader := bufio.NewReader(os.Stdin)
	// fmt.Print("www.energiatili.fi username: ")
	// username, err := reader.ReadString('\n')
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println("www.energiatili.fi password:")
	// password, err := terminal.ReadPassword(0)
	// if err != nil {
	// 	panic(err)
	// }

	f, err := fetcher.New()
	if err != nil {
		panic(err)
	}

	fmt.Println("Logging in…")
	// err = f.Login(username, string(password))
	err = f.Login(secrets.Username(), secrets.Password())
	if err != nil {
		panic(err)
	}

	fmt.Println("Login OK. Downloading consumption data…")
	if err := f.ConsumptionReport(os.Stdout); err != nil {
		panic(err)
	}
}
