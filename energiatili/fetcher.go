// Package energiatili downloads www.energiatili.fi energy consumption JSON data
package energiatili

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

// ErrorStatusCode is returned if server returns an unexpected HTTP status code
var ErrorStatusCode = errors.New("unexpected HTTP status code from server")

// Client retrieves data from Energiatili
type Client struct {
	// LoginURL is address for the login request to initialize session
	LoginURL string

	// ConsumptionReportURL consumption data HTML download URL
	ConsumptionReportURL string

	// GetUsernamePassword is called to get credentials to log in to service
	GetUsernamePassword func() (string, string, error)

	// unexported fields
	cl       *http.Client
	loggedIn bool
}

// ConsumptionReport fetches the actual consumption report data (JSON)
func (f *Client) ConsumptionReport(w io.Writer) (err error) {
	f.initialize()
	if !f.loggedIn {
		if err := f.login(); err != nil {
			return err
		}
	}
	resp, err := f.cl.Get(f.ConsumptionReportURL)
	if err != nil {
		return
	}
	if resp.StatusCode != http.StatusOK {
		return ErrorStatusCode
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	data, err := html2json(body)
	if err != nil {
		return err
	}
	r := strings.NewReader(data)
	_, err = io.Copy(w, r)
	return err
}

// New initializes a fetcher with a fresh cookie jar
func (f *Client) initialize() {
	if f.cl == nil {
		jar, _ := cookiejar.New(nil)
		f.cl = &http.Client{Jar: jar}
	}
}

func (f *Client) login() (err error) {
	f.initialize()
	username, password, err := f.GetUsernamePassword()
	if err != nil {
		return
	}
	form := url.Values{
		"username": []string{username},
		"password": []string{password},
	}
	resp, err := f.cl.PostForm(f.LoginURL, form)
	if err != nil {
		return
	}
	if resp.StatusCode != http.StatusOK {
		return ErrorStatusCode
	}
	f.loggedIn = true
	return
}
