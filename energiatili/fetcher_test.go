package energiatili

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// var runIntegrationTests = flag.Bool("integration", false, "Run the integration tests (in addition to the unit tests)")

const (
	testLoginUser     = "testUserName"
	testLoginPassword = "p4ssw0rdForTest"
)

// TestLoginStatus tests error handling when server returns 403 Forbidden
func TestLoginStatus(t *testing.T) {
	ts := testServer{}
	ts.Start()
	ts.statusCode = 403
	defer ts.Close()

	fetcher := Client{
		GetUsernamePassword:  mockGetUsernamePassword,
		LoginURL:             ts.URL,
		ConsumptionReportURL: ts.URL,
	}
	err := fetcher.login()
	if err == nil {
		t.Error("Login did not return error; expected error when HTTP status 403")
	}

}

// TestNoResponse tests error handling when server closes connection
func TestNoResponse(t *testing.T) {
	ts := testServer{}
	ts.Start()
	ts.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ts.CloseClientConnections()
	})
	defer ts.Close()

	fetcher := Client{
		GetUsernamePassword:  mockGetUsernamePassword,
		LoginURL:             ts.URL,
		ConsumptionReportURL: ts.URL,
	}
	err := fetcher.login()
	if err == nil {
		t.Error("Login did not return error; expected error when HTTP status 403")
	}

}

// TestLoginForm tests the username and password are sent to the test server
func TestLoginForm(t *testing.T) {
	ts := testServer{}
	ts.Start()
	defer ts.Close()

	fetcher := Client{
		GetUsernamePassword:  mockGetUsernamePassword,
		LoginURL:             ts.URL,
		ConsumptionReportURL: ts.URL,
	}
	fetcher.login()
	if len(ts.requests) != 1 {
		t.Errorf("Expected requests count=1, got count=%v", len(ts.requests))
	}
	req := ts.requests[len(ts.requests)-1]
	expectedForm := map[string]string{
		"username": testLoginUser,
		"password": testLoginPassword,
	}
	for key, want := range expectedForm {
		if req.FormValue(key) != want {
			t.Errorf("Expected request form %v=%v, got %v", key, want, req.FormValue(key))
		}
	}
}

// TestConsumptionReport tests the full flow
func TestConsumptionReport(t *testing.T) {
	ts := testServer{}
	ts.Start()
	defer ts.Close()

	fetcher := Client{
		GetUsernamePassword:  mockGetUsernamePassword,
		LoginURL:             ts.URL,
		ConsumptionReportURL: ts.URL,
	}
	err := fetcher.login()
	if err != nil {
		t.Errorf("Got unexpected error from Login(): %v", err)
	}
	// Mock response that looks like the real thing
	ts.body = `<html>
<p>Some random stuff here</p>
Then magically, var model = {"first": "value", "second": new Date(1234)} ;
More stuff
</html>`
	buf := &bytes.Buffer{}
	err = fetcher.ConsumptionReport(buf)
	if err != nil {
		t.Errorf("Got unexpected error from ConsumptionReport(): %v", err)
	}
	want := `{"first": "value", "second": new Date(1234)} `
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

	*httptest.Server
}

func (t *testServer) Start() {
	t.Server = httptest.NewServer(t)
	t.statusCode = 200
	t.body = "fetcher_test Unit Test server response body"
}

// HTTP 200 ok with simple text body
func (t *testServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Set-Cookie", ".ASPXAUTH=test_auth_value")
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(t.statusCode)
	fmt.Fprintf(w, t.body)
	_ = r.ParseForm()
	t.requests = append(t.requests, *r)
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

func mockGetUsernamePassword() (username, password string, err error) {
	return testLoginUser, testLoginPassword, nil
}
