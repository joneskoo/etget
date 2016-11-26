// Package energiatili downloads www.energiatili.fi energy consumption JSON data
package energiatili

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

var errorStatusCode = errors.New("unexpected status code from server")

// authCookieKey must be set in login or login is considered failed
const authCookieKey = ".ASPXAUTH"

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
func (f *Client) ConsumptionReport(w io.Writer) error {
	if !f.loggedIn {
		if err := f.login(); err != nil {
			// Login again and try again
			if err := f.login(); err != nil {
				return err
			}
		}
	}
	bodyBytes, err := f.requestConsumptionReport()
	if err != nil {
		return err
	}

	// Find var model = ....
	startData := "var model = "
	endData := ";"
	body := string(bodyBytes)
	start := strings.Index(body, startData)
	if start == -1 {
		return fmt.Errorf("failed to find %q in body", startData)
	}
	end := start + strings.Index(body[start:], endData)
	if end == -1 {
		return fmt.Errorf("unterminated %q in body", startData)
	}
	fmt.Fprint(w, body[start+len(startData):end])
	return nil
}

func (f *Client) requestConsumptionReport() (body []byte, err error) {
	resp, err := f.cl.Get(f.ConsumptionReportURL)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		io.Copy(ioutil.Discard, resp.Body)
		return nil, fmt.Errorf("unexpected status code %d from server", resp.StatusCode)
	}
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (f *Client) login() (err error) {
	if f.cl == nil {
		jar, _ := cookiejar.New(nil)
		f.cl = &http.Client{Jar: jar}
	}
	username, password, err := f.GetUsernamePassword()
	if err != nil {
		return err
	}
	form := url.Values{
		"username": []string{username},
		"password": []string{password},
	}
	resp, err := f.cl.PostForm(f.LoginURL, form)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	io.Copy(ioutil.Discard, resp.Body)

	url, _ := url.Parse(f.LoginURL)
	for _, cookie := range f.cl.Jar.Cookies(url) {
		if cookie.Name == authCookieKey {
			f.loggedIn = true
		}
	}
	if !f.loggedIn {
		return fmt.Errorf("login did not return authentication cookie %s", authCookieKey)
	}

	if resp.StatusCode != http.StatusOK {
		return errorStatusCode
	}
	return nil
}
