package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/joneskoo/etget/fetcher"
)

func main() {
	data, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}
	json, err := fetcher.HTML2JSON(data)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v\n", string(json))
}
