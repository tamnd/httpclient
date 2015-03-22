package httpclient

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
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
type Client struct {
	client *http.Client
}

// New returns new client.
func New(client *http.Client) *Client {
	return &Client{client: client}
}

func (c *Client) err(resp *http.Response, message string) error {
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
func (c *Client) Get(url string) (*http.Response, error) {
	return c.client.Get(url)
}

// Bytes fetches the specified url and returns the response body as bytes.
func (c *Client) Bytes(url string) ([]byte, error) {
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
func (c *Client) String(url string) (string, error) {
	bytes, err := c.Bytes(url)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// Reader issues a GET request to a specified URL and returns an reader from the response body.
func (c *Client) Reader(url string) (io.ReadCloser, error) {
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
func (c *Client) JSON(url string, v interface{}) error {
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
func (c *Client) XML(url string, v interface{}) error {
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
func (c *Client) Files(urls []string, files []*File) error {
	ch := make(chan error, len(files))
	for i := range files {
		go func(i int) {
			resp, err := c.Get(urls[i])
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
			files[i].Data, err = ioutil.ReadAll(resp.Body)
			if err != nil {
				ch <- c.err(resp, err.Error())
				return
			}
			ch <- nil
		}(i)
	}
	for range files {
		if err := <-ch; err != nil {
			return err
		}
	}
	return nil
}

// Download downloads multiple files concurrency.
func (c *Client) Download(urls []string, files []*File) error {
	return c.Files(urls, files)
}

var client = &Client{client: &http.Client{}}

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
func Files(urls []string, files []*File) error {
	return client.Files(urls, files)
}

// Download downloads multiple files concurrency.
func Download(urls []string, files []*File) error {
	return client.Files(urls, files)
}
