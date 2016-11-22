// Package energiatili downloads www.energiatili.fi energy consumption JSON data
package energiatili

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
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
	_, err = html2json(string(body), w)
	return err
}

func (f *Client) login() (err error) {
	if f.cl == nil {
		jar, _ := cookiejar.New(nil)
		f.cl = &http.Client{Jar: jar}
	}
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
