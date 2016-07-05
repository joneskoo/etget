// Package fetcher downloads www.energiatili.fi energy consumption JSON data
package fetcher

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

var (
	// LoginURL is login request URL
	LoginURL = "https://www.energiatili.fi/Extranet/Extranet/LogIn"
	// ConsumptionReportURL consumption data HTML download URL
	ConsumptionReportURL = "https://www.energiatili.fi/Reporting/CustomerConsumption/UserConsumptionReport"
)

// Fetcher holds the login state such as cookies
type Fetcher struct {
	// unexported fields
	cl       http.Client
	loggedIn bool
}

// ErrorNotLoggedIn is raised if trying to call url fetch before Login()
var (
	ErrorNotLoggedIn = errors.New("NeedLoginFirst")
	ErrorStatusCode  = errors.New("UnexpectedHTTPStatusCodeFromServer")
)

// New initializes a fetcher with a fresh cookie jar
func New() (f Fetcher, err error) {
	jar, _ := cookiejar.New(nil)
	f = Fetcher{
		cl: http.Client{
			Jar: jar,
		},
	}
	return
}

// Login logs in to www.energiatili.fi
func (f *Fetcher) Login(user, password string) (err error) {
	form := url.Values{
		"username": []string{user},
		"password": []string{password},
	}
	resp, err := f.cl.PostForm(LoginURL, form)
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
	if f.loggedIn == false {
		return ErrorNotLoggedIn
	}
	resp, err := f.cl.Get(ConsumptionReportURL)
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
