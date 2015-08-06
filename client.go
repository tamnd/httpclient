package httpclient

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sync"
)

// Error is the custom error type returns from HTTP requests.
type Error struct {
	Message    string
	StatusCode int
	URL        string
}

// Error returns the error message.
func (e *Error) Error() string {
	return e.Message
}

// File represents a file.
type File struct {
	// File name with no directory.
	Name string

	// Contents of the file.
	Data []byte
}

// A Client is an HTTP client.
// It wraps net/http's client and add some methods for making HTTP request easier.
type httpClient struct {
	client *http.Client
}

// New returns new client.
func New() *httpClient {
	return &httpClient{client: &http.Client{}}
}

func (c *httpClient) err(resp *http.Response, message string) error {
	if message == "" {
		message = fmt.Sprintf("Get %s -> %d", resp.Request.URL.String(), resp.StatusCode)
	}
	return &Error{
		Message:    message,
		StatusCode: resp.StatusCode,
		URL:        resp.Request.URL.String(),
	}
}

// Get issues a GET to the specified URL. It returns an http.Response for further processing.
func (c *httpClient) Get(url string) (*http.Response, error) {
	return c.client.Get(url)
}

// Bytes fetches the specified url and returns the response body as bytes.
func (c *httpClient) Bytes(url string) ([]byte, error) {
	resp, err := c.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, c.err(resp, "")
	}
	p, err := ioutil.ReadAll(resp.Body)
	return p, err
}

// String fetches the specified URL and returns the response body as a string.
func (c *httpClient) String(url string) (string, error) {
	bytes, err := c.Bytes(url)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// Reader issues a GET request to a specified URL and returns an reader from the response body.
func (c *httpClient) Reader(url string) (io.ReadCloser, error) {
	resp, err := c.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		err = c.err(resp, "")
		resp.Body.Close()
		return nil, err
	}
	return resp.Body, nil
}

// JSON issues a GET request to a specified URL and unmarshal json data from the response body.
func (c *httpClient) JSON(url string, v interface{}) error {
	resp, err := c.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return c.err(resp, "")
	}
	err = json.NewDecoder(resp.Body).Decode(v)
	if _, ok := err.(*json.SyntaxError); ok {
		err = c.err(resp, "JSON syntax error at "+url)
	}
	return err
}

// XML issues a GET request to a specified URL and unmarshal XML data from the response body.
func (c *httpClient) XML(url string, v interface{}) error {
	resp, err := c.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return c.err(resp, "")
	}
	err = xml.NewDecoder(resp.Body).Decode(v)
	return err
}

// Files downloads multiple files concurrency.
func (c *httpClient) Files(urls []string, files *[]File) error {
	l := len(urls)
	fs := make([]File, l)
	ch := make(chan error, l)
	var wg sync.WaitGroup
	wg.Add(l)
	for i, url := range urls {
		go func(i int) {
			defer wg.Done()
			resp, err := c.Get(url)
			if err != nil {
				ch <- err
				return
			}
			defer resp.Body.Close()
			if resp.StatusCode != 200 {
				var err error
				err = c.err(resp, "")
				ch <- err
				return
			}
			fs[i].Data, err = ioutil.ReadAll(resp.Body)
			if err != nil {
				ch <- c.err(resp, err.Error())
				return
			}
			ch <- nil
		}(i)
	}
	wg.Wait()
	for _ = range fs {
		if err := <-ch; err != nil {
			return err
		}
	}
	*files = fs
	return nil
}

// Download downloads multiple files concurrency.
func (c *httpClient) Download(urls []string, files *[]File) error {
	return c.Files(urls, files)
}

var client = New()

// Get issues a GET to the specified URL. It returns an http.Response for further processing.
func Get(url string) (*http.Response, error) {
	return client.Get(url)
}

// Bytes fetches the specified url and returns the response body as bytes.
func Bytes(url string) ([]byte, error) {
	return client.Bytes(url)
}

// String fetches the specified URL and returns the response body as a string.
func String(url string) (string, error) {
	return client.String(url)
}

// Reader issues a GET request to a specified URL and returns an reader from the response body.
func Reader(url string) (io.ReadCloser, error) {
	return client.Reader(url)
}

// JSON issues a GET request to a specified URL and unmarshal json data from the response body.
func JSON(url string, v interface{}) error {
	return client.JSON(url, v)
}

// XML issues a GET request to a specified URL and unmarshal xml data from the response body.
func XML(url string, v interface{}) error {
	return client.JSON(url, v)
}

// Files downloads multiple files concurrency.
func Files(urls []string, files *[]File) error {
	return client.Files(urls, files)
}

// Download downloads multiple files concurrency.
func Download(urls []string, files *[]File) error {
	return client.Files(urls, files)
}
