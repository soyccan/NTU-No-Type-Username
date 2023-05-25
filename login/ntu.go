package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/http/httputil"
	"net/url"
	"regexp"
	"strings"
)

const LogHTTPMsg = false

// credentials struct
type credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (c *credentials) valid() bool {
	return len(c.Username) > 0
}

func loadCredentials(path string) *credentials {
	var cred credentials

	fileData, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal("Error reading JSON file:", err)
		return nil
	}

	err = json.Unmarshal(fileData, &cred)
	if err != nil || !cred.valid() {
		log.Fatal("Error parsing JSON data:", err)
		return nil
	}
	return &cred
}

// implements http.RoundTripper
type myTransport struct{}

// hook before sending data
// log all http requests for debugging
func (*myTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Workaround: the web server rejects port number in Host field of HTTP header
	// so we tamper url host by removing default port number
	host, port, cut := strings.Cut(req.URL.Host, ":")
	if cut && (port == "80" || port == "443") {
		req.URL.Host = host
	}

	if !LogHTTPMsg {
		resp, err := http.DefaultTransport.RoundTrip(req)
		if err != nil {
			log.Fatal(err)
		}
		return resp, err
	}

	reqData, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		log.Fatal(err)
	}

	resp, rtErr := http.DefaultTransport.RoundTrip(req)
	if rtErr != nil {
		log.Fatal(rtErr)
	}

	respData, err := httputil.DumpResponse(resp, true)
	if err != nil {
		log.Fatal(err)
	}
	reqData = append(reqData, respData...)

	fmt.Printf("%s\n", reqData)

	return resp, rtErr
}

// main crawler struct
type NoTypeUsername struct {
	CredPath   string
	SamlPath   string
	client     *http.Client
	noRedirect bool
}

func (ntu *NoTypeUsername) initHTTPClient() {
	if ntu.client != nil {
		return
	}
	ntu.client = &http.Client{
		Jar: func(j *cookiejar.Jar, _ error) *cookiejar.Jar { return j }(cookiejar.New(nil)),
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if ntu.noRedirect {
				// Disable automatic redirect by returning http.ErrUseLastResponse
				ntu.noRedirect = false
				return http.ErrUseLastResponse
			}
			return nil
		},
		Transport: &myTransport{},
	}
}

// get NTU Single-Sign-On URL by accessing NTU COOL and redirecting
func (ntu *NoTypeUsername) getSAMLURL(loginURL string) string {
	ntu.initHTTPClient()

	// Send GET request to NTU COOL login page
	log.Println("Visiting login page for the first time...")
	ntu.noRedirect = true
	resp, err := ntu.client.Get(loginURL)
	if err != nil {
		log.Fatalf("Failed visiting login page %s: %s\n", loginURL, err)
		return ""
	}
	defer resp.Body.Close()

	// Get redirection location
	loc := resp.Header["Location"]
	if len(loc) == 0 {
		log.Fatal("Failed to get SSO URL")
		return ""
	}
	return loc[0]
}

func findSAMLResponse(html []byte) string {
	re := regexp.MustCompile(`name="SAMLResponse" value="([^"]+)"`)
	match := re.FindSubmatch(html)
	if len(match) < 2 {
		return ""
	}
	return string(match[1])
}

func (ntu *NoTypeUsername) authSAML(loginURL string) string {
	ntu.initHTTPClient()

	samlURL := ntu.getSAMLURL(loginURL)
	if samlURL == "" {
		return ""
	}

	cred := loadCredentials(ntu.CredPath)
	if cred == nil {
		return ""
	}

	formData := url.Values{
		"__VIEWSTATE":          {"/wEPDwUKMTY2MTc3NjUzM2RkUK4S8IU/lZeKUDrQIAtt4tRhRV4ZOkEMNdoJavm/SBs="},
		"__VIEWSTATEGENERATOR": {"0EE29E36"},
		"__EVENTVALIDATION":    {"/wEdAAUdVdOEjcCKz7S6sLphMAmFlt/S8mKmQpmuxn2LW6B9thvLC/FQOf5u4GfePSXQdrRBPkcB0cPQF9vyGTuIFWmijKZWG4rH59f66Vc64WGnN/Hmf00Q2eMalQURbQ6cPb45rGUVCHnIwpyxWjkkPDce"},
		"__db":                 {"15"},
		"ctl00$ContentPlaceHolder1$UsernameTextBox": {cred.Username},
		"ctl00$ContentPlaceHolder1$PasswordTextBox": {cred.Password},
		"ctl00$ContentPlaceHolder1$SubmitButton":    {"Sign+In"},
	}
	log.Println("Login SAML...")
	resp, err := ntu.client.PostForm(samlURL, formData)
	if err != nil {
		log.Fatalf("POST SAML %s failed: %s\n", samlURL, err)
		return ""
	}
	defer resp.Body.Close()

	// Extract SAML response from the response body
	html, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading response body:", err)
		return ""
	}
	return findSAMLResponse(html)
}

func (ntu *NoTypeUsername) LoginNTUCOOL() {
	coolURL := "https://cool.ntu.edu.tw/login/saml"
	samlResp := ntu.authSAML(coolURL)
	if samlResp == "" {
		log.Fatal("SAMLResponse not found")
		return
	}

	formData := url.Values{"SAMLResponse": {samlResp}}
	resp, err := ntu.client.PostForm(coolURL, formData)
	if err != nil {
		log.Fatalf("Visit %s error: %s\n", coolURL, err)
		return
	}

	html, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading response body:", err)
		return
	}
	fmt.Println(string(html))

	// write SAMLResponse to file
	err = ioutil.WriteFile(ntu.SamlPath, []byte(samlResp), 0666)
	if err != nil {
		log.Fatal("Write to SAML file fail:", err)
	}
}
