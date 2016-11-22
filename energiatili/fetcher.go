// Package fetcher downloads www.energiatili.fi energy consumption JSON data
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

// Fetcher holds the login state such as cookies
type Fetcher struct {
	// LoginURL is login request URL
	LoginURL string

	// ConsumptionReportURL consumption data HTML download URL
	ConsumptionReportURL string

	// GetUsernamePassword is called to get credentials to log in to service
	GetUsernamePassword func() (string, string, error)

	// unexported fields
	cl       *http.Client
	loggedIn bool
}

// ErrorNotLoggedIn is raised if trying to call url fetch before Login()
var (
	ErrorNotLoggedIn = errors.New("NeedLoginFirst")
	ErrorStatusCode  = errors.New("UnexpectedHTTPStatusCodeFromServer")
)

// New initializes a fetcher with a fresh cookie jar
func (f *Fetcher) initialize() {
	if f.cl == nil {
		jar, _ := cookiejar.New(nil)
		f.cl = &http.Client{Jar: jar}
	}
}

// Login logs in to www.energiatili.fi
func (f *Fetcher) Login(user, password string) (err error) {
	f.initialize()
	form := url.Values{
		"username": []string{user},
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

// ConsumptionReport fetches the actual consumption report data (JSON)
func (f *Fetcher) ConsumptionReport(w io.Writer) (err error) {
	f.initialize()
	if f.loggedIn == false {
		return ErrorNotLoggedIn
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
