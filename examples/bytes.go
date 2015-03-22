package main

import (
	"fmt"
	"log"

	"github.com/tamnd/httpclient"
)

func main() {
	content, err := httpclient.Bytes("http://www.example.com")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("%#v", content)
}
