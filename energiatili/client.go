// Package energiatili downloads www.energiatili.fi energy consumption JSON data
package energiatili

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

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

var newDateReplacer = strings.NewReplacer("new Date(", "", ")", "")

// ConsumptionReport fetches the actual consumption report data (JSON)
func (c *Client) ConsumptionReport(w io.Writer) error {
	if !c.loggedIn {
		if err := c.login(); err != nil {
			// Login again and try again
			if err := c.login(); err != nil {
				return err
			}
		}
	}
	bodyBytes, err := c.requestConsumptionReport()
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
	newDateReplacer.WriteString(w, body[start+len(startData):end])
	return nil
}

func (c *Client) requestConsumptionReport() (body []byte, err error) {
	resp, err := c.cl.Get(c.ConsumptionReportURL)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		io.Copy(ioutil.Discard, resp.Body)
		return nil, fmt.Errorf("want HTTP status code 200, got %d", resp.StatusCode)
	}
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (c *Client) login() (err error) {
	if c.cl == nil {
		jar, _ := cookiejar.New(nil)
		c.cl = &http.Client{Jar: jar}
	}
	username, password, err := c.GetUsernamePassword()
	if err != nil {
		return err
	}
	form := url.Values{
		"username": []string{username},
		"password": []string{password},
	}
	resp, err := c.cl.PostForm(c.LoginURL, form)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	io.Copy(ioutil.Discard, resp.Body)

	url, _ := url.Parse(c.LoginURL)
	for _, cookie := range c.cl.Jar.Cookies(url) {
		if cookie.Name == authCookieKey {
			c.loggedIn = true
		}
	}
	if !c.loggedIn {
		return fmt.Errorf("login did not set authentication cookie %s", authCookieKey)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("want HTTP status code 200, got %d", resp.StatusCode)

	}
	return nil
}
