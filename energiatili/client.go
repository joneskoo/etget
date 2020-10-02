// Package energiatili downloads www.energiatili.fi energy consumption JSON data
package energiatili

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"sync"
)

const (
	endpointLogin             = "https://www.energiatili.fi/Extranet/Extranet/LogIn"
	endpointConsumptionReport = "https://www.energiatili.fi/Reporting/CustomerConsumption/UserConsumptionReport"
	endpointFullConsumption   = "https://www.energiatili.fi/Reporting/SessionlessConsumption/SessionlessGetFullConsumption"
)

// Client retrieves data from Energiatili.
type Client struct {
	// UsernamePasswordFunc is called on login to acquire credentials.
	UsernamePasswordFunc func() (username string, password string, err error)

	// Transport is a roundtripper the client uses to make HTTP requests.
	Transport http.RoundTripper

	// unexported
	initOnce sync.Once
	cl       http.Client
}

func (c *Client) init() {
	c.initOnce.Do(func() {
		jar, _ := cookiejar.New(nil)

		c.cl = http.Client{
			Transport: c.Transport,
			Jar:       jar,
		}
	})
}

// ConsumptionReport fetches the actual consumption report data (JSON)
func (c *Client) ConsumptionReport(ctx context.Context, w io.Writer) error {
	c.init()

	if err := c.login(ctx); err != nil {
		return err
	}

	// Dummy request - without this the real request will return empty.
	req, err := http.NewRequestWithContext(ctx, "GET", endpointConsumptionReport, nil)
	if err != nil {
		return err
	}
	resp, err := c.cl.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		io.Copy(ioutil.Discard, resp.Body)
		return fmt.Errorf("want HTTP status code 200, got %d", resp.StatusCode)
	}
	bodyBytes, _ := ioutil.ReadAll(resp.Body)

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

var newDateReplacer = strings.NewReplacer("new Date(", "", ")", "")

func (c *Client) login(ctx context.Context) (err error) {
	username, password, err := c.UsernamePasswordFunc()
	if err != nil {
		return err
	}

	form := url.Values{
		"username": []string{username},
		"password": []string{password},
	}
	req, err := http.NewRequestWithContext(ctx, "POST", endpointLogin, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.cl.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	io.Copy(ioutil.Discard, resp.Body)

	u, err := url.Parse(endpointLogin)
	if err != nil {
		panic(err)
	}
	for _, cookie := range c.cl.Jar.Cookies(u) {
		if cookie.Name == ".ASPXAUTH" {
			return nil
		}
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("want HTTP status code 200, got %d", resp.StatusCode)
	}
	return fmt.Errorf("login did not set expected authentication cookie .ASPXAUTH")
}
