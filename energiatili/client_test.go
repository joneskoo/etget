package energiatili_test

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/joneskoo/etget/energiatili"
)

// var runIntegrationTests = flag.Bool("integration", false, "Run the integration tests (in addition to the unit tests)")

const (
	testLoginUser     = "testUserName"
	testLoginPassword = "p4ssw0rdForTest"
)

// TestLoginStatus tests error handling when server returns 403 Forbidden
func TestLoginStatus(t *testing.T) {
	ts := &testServer{
		statusCode: 403,
		body:       "fetcher_test Unit Test server response body",
	}

	fetcher := energiatili.Client{
		UsernamePasswordFunc: mockUsernamePasswordFunc,
		Transport:            ts,
	}
	err := fetcher.ConsumptionReport(context.TODO(), ioutil.Discard)
	if err == nil {
		t.Error("Login did not return error; expected error when HTTP status 403")
	}

}

// TestLoginForm tests the username and password are sent to the test server
func TestLoginForm(t *testing.T) {
	ts := &testServer{
		statusCode: 200,
		body:       "fetcher_test Unit Test server response body",
	}

	fetcher := energiatili.Client{
		UsernamePasswordFunc: mockUsernamePasswordFunc,
		Transport:            ts,
	}
	fetcher.ConsumptionReport(context.TODO(), ioutil.Discard)
	if len(ts.requests) != 2 {
		t.Errorf("want 2 requests, got count=%d", len(ts.requests))
	}
	req := ts.requests[0] // login
	expectedForm := map[string]string{
		"username": testLoginUser,
		"password": testLoginPassword,
	}
	for key, want := range expectedForm {
		if req.FormValue(key) != want {
			t.Errorf("want format %v=%v, got %q", key, want, req.FormValue(key))
		}
	}
}

// TestConsumptionReport tests the full flow
func TestConsumptionReport(t *testing.T) {
	ts := &testServer{
		statusCode: 200,
		body:       "fetcher_test Unit Test server response body",
	}

	fetcher := energiatili.Client{
		UsernamePasswordFunc: mockUsernamePasswordFunc,
		Transport:            ts,
	}
	// Mock response that looks like the real thing
	ts.body = `<html>
<p>Some random stuff here</p>
Then magically, var model = {"first": "value", "second": new Date(1234)} ;
More stuff
</html>`
	buf := &bytes.Buffer{}
	err := fetcher.ConsumptionReport(context.TODO(), buf)
	if err != nil {
		t.Errorf("Got unexpected error from ConsumptionReport(): %v", err)
	}
	want := `{"first": "value", "second": 1234} `
	if buf.String() != want {
		t.Errorf("ConsumptionReport(w) want %q written, got %q", want, buf.String())
		t.FailNow()
	}
}

type testServer struct {
	oldLoginURL       string
	oldConsumptionURL string
	handler           http.Handler
	statusCode        int
	body              string
	requests          []http.Request
}

// HTTP 200 ok with simple text body
func (t *testServer) RoundTrip(req *http.Request) (res *http.Response, err error) {
	res = &http.Response{
		StatusCode: t.statusCode,
		Body:       ioutil.NopCloser(strings.NewReader(t.body)),
		Header:     make(http.Header),
	}
	res.Header.Set("Set-Cookie", ".ASPXAUTH=test_auth_value")
	res.Header.Set("Content-Type", "text/plain; charset=utf-8")
	res.Header.Set("X-Content-Type-Options", "nosniff")
	_ = req.ParseForm()
	t.requests = append(t.requests, *req)
	return res, nil
}

// // TestFetcherIntegrationTests tests login and fetching data against the real service
// func TestFetcherIntegrationTests(t *testing.T) {
// 	if !*runIntegrationTests {
// 		t.Skip("To run this test, use: go test -integration")
// 	}
// 	fetcher, _ := New()
// 	err := fetcher.Login(testLoginUser, testLoginPassword)
// 	if err != nil {
// 		t.Error(err)
// 		t.Fail()
// 	}
// 	_, err = fetcher.ConsumptionReport()
// 	if err != nil {
// 		t.Errorf("Expected ConsumptionReport() to succeed, but got error: %v", err)
// 	}
// }

func mockUsernamePasswordFunc() (username, password string, err error) {
	return testLoginUser, testLoginPassword, nil
}
