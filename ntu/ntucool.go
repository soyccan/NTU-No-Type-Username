package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

type NtuCOOL struct {
	NoTypeUsername
}

// get NTU Single-Sign-On URL by accessing NTU COOL and redirecting
func (ntu *NtuCOOL) getSAMLURL(loginURL string) (samlURL string, err error) {
	err = ntu.initHTTPClient()
	if err != nil {
		return
	}

	// Send GET request to NTU COOL login page
	log.Println("Visiting login page for the first time...")
	ntu.noRedirect = true
	resp, err := ntu.client.Get(loginURL)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	// Get redirection location
	loc := resp.Header["Location"]
	if len(loc) == 0 {
		return
	}
	samlURL = loc[0]
	return
}

func findSAMLResponse(html []byte) string {
	re := regexp.MustCompile(`name="SAMLResponse" value="([^"]+)"`)
	match := re.FindSubmatch(html)
	if len(match) < 2 {
		return ""
	}
	return string(match[1])
}

func (ntu *NtuCOOL) authSAML(loginURL string) (samlResp string, err error) {
	err = ntu.initHTTPClient()
	if err != nil {
		return
	}

	samlURL, err := ntu.getSAMLURL(loginURL)
	if err != nil {
		return
	}

	cred, err := loadCredentials(ntu.CredPath)
	if err != nil {
		return
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
		return
	}
	defer resp.Body.Close()

	// Extract SAML response from the response body
	html, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading response body:", err)
		return
	}
	samlResp = findSAMLResponse(html)
	return
}

func (ntu *NtuCOOL) checkSucceedLogin(coolURL, samlResp string) (
	succeed bool, resp *http.Response, err error,
) {
	formData := url.Values{"SAMLResponse": {samlResp}}
	resp, err = ntu.client.PostForm(coolURL, formData)
	if err != nil {
		log.Fatalf("Visit %s error: %s\n", coolURL, err)
		return
	}

	succeed = (resp.StatusCode == 200 &&
		resp.Request.URL.Query().Get("login_success") == "1")
	return
}

// Convert cookies into HTTP header format
// Example: "Cookies: log_session_id=c550c10407e11107b156f34ffbd60c47; _csrf_token=gKU%2Fci5TKNlELoFuUqK2FWY6tON4vvDwxx"
func formatCookiesHeader(jar http.CookieJar, url_ string) string {
	var cookies []string
	u, err := url.Parse(url_)
	if err != nil {
		log.Fatal("url parse failed")
	}
	for _, c := range jar.Cookies(u) {
		cookies = append(cookies, c.String())
	}
	return strings.Join(cookies, "; ")
}

// Convert cookies into a NGINX config that set the cookies to a var
// Example:
func formatCookiesNginx(jar http.CookieJar, url_ string) string {
	return fmt.Sprintf("set $cookies \"%s\";", formatCookiesHeader(jar, url_))
}

func (ntu *NtuCOOL) Login() (err error) {
	coolURL := "https://cool.ntu.edu.tw/login/saml"
	samlResp, err := ntu.authSAML(coolURL)
	if samlResp == "" {
		log.Fatal("SAMLResponse not found")
		return
	}

	succeed, resp, err := ntu.checkSucceedLogin(coolURL, samlResp)
	if !succeed || err != nil {
		var html []byte
		html, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal("Error reading response body:", err)
			return
		}
		log.Println(string(html))
		log.Fatal("Login COOL failed")
		return
	}

	// write cookies to file
	err = ioutil.WriteFile(
		ntu.CookiePath,
		[]byte(formatCookiesNginx(ntu.client.Jar, coolURL)),
		0666,
	)
	if err != nil {
		log.Fatal("Write to SAML file fail:", err)
		return
	}
	return
}
