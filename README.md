# HttpClient

[![Software License](https://img.shields.io/badge/license-MIT-brightgreen.svg?style=flat-square)](LICENSE.md)
[![Build Status](https://img.shields.io/travis/tamnd/httpclient/master.svg?style=flat-square)](https://travis-ci.org/tamnd/httpclient)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/tamnd/httpclient)

Simplified HTTP Client for [Go](http://www.golang.org). 

HttpClient is designed to be the simplest way possible to make http requests. It sends an HTTP request and unmarshals json, xml, csv from the response in just one function call.

```go
package main

import (
	"fmt"
	"log"

	"github.com/tamnd/httpclient"
)

type Person struct {
	ID        string
	FirstName string
	Gender    string
	LastName  string
	Link      string
	Locate    string
	Name      string
	Username  string
}

func main() {
	var mark Person
	err := httpclient.JSON("http://graph.facebook.com/4", &mark)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%#v\n", mark)
}
```

Output:

```
main.Person{ID:"4", FirstName:"", Gender:"male", LastName:"", Link:"https://www.facebook.com/zuck", Locate:"", Name:"Mark Zuckerberg", Username:"zuck"}
```

## Features
- Get and unmarshal JSON from a url
- Get and unmarshal XML from a url
- Download multipe files concurrency

## Install
```
$ go get github.com/tamnd/httpclient
```

## Usage

```go
import "tamnd/misc/httpclient"
```

## How to
- [Get String](#get-string) 
- [Get Bytes](#get-bytes) 
- [Get JSON](#get-json) 
- [Get XML](#get-xml)
- [Get Reader](#get-reader)
- [Download Files](#download-files)
- Send POST Request
- Upload Files
- Custom Request Header

### Get String

```go
func String(url string) (string, error)
```
String fetches the specified URL and returns the response body as a string.

```go
content, err := http.Client.String("http://www.example.com")
```

### Get Bytes

```go
func Bytes(url string) ([]byte, error)
```
Bytes fetches the specified url and returns the response body as bytes.

```go
bytes, err := http.Client.Bytes("http://www.example.com")
```

### Get JSON

```go
func JSON(url string, v interface{}) error
```
JSON issues a GET request to a specified URL and unmarshal json data from the
response body.

```go
var user = struct {
    ID        string `json:"id"`
    Gender    string `json:"gender"`
    Link      string `json:"link"`
    Name      string `json:"name"`
    Username  string `json:"username"`
}{}
err := httpclient.JSON("http://graph.facebook.com/4", &user)
```

### Download Files

```go
func Download(urls []string, files []*File) error
```
Download downloads multiple files concurrency.

```go
urls := []string{
	"http://www.golang.org",
	"http://www.clojure.org",
	"http://www.erlang.org",
}
var files []*httpclient.File
err := httpclient.Download(urls, files)
```


#### Get Reader from Response

```go
func Reader(url string) (io.ReadCloser, error)
```
Reader issues a GET request to a specified URL and returns an reader from the
response body.


## Roadmap
- [ ] Send POST request
- [ ] Custom request header
- [ ] Send basic authentication
- [ ] Make `Upload()` function
- [ ] Get response header
- [ ] Connection timeoutspackage main
- [ ] Custom error handling

## Contribute

- Fork repository
- Create a feature branch
- Open a new pull request
- Create an issue for bug report or feature request

## Contact

- Nguyen Duc Tam
- [tamnd87@gmail.com](mailto:tamnd87@gmail.com)
- [http://twitter.com/tamnd87](http://twitter.com/tamnd87)

## License
The MIT License (MIT). Please see [LICENSE](LICENSE) for more information.

Copyright (c) 2015 Nguyen Duc Tam, tamnd87@gmail.com
