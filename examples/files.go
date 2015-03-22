package main

import (
	"fmt"
	"log"

	"github.com/tamnd/httpclient"
)

func main() {
	urls := []string{
		"http://www.golang.org",
		"http://www.clojure.org",
		"http://www.erlang.org",
	}
	var files []*httpclient.File
	err := httpclient.Download(urls, files)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Download files success!")
}
