package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"text/template"
)

func RequestURL(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("home.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func RedirectHandler(w http.ResponseWriter, r *http.Request) {
	m, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	code := m.Get("code")
	scope := m.Get("scope")
	fmt.Fprintln(w, "Code: ", code, " Scope: ", scope)
	accessToken = getAccessToken(w, r, code)

}

func GetUserData(w http.ResponseWriter, r *http.Request) {
	if accessToken == "" {
		fmt.Fprintf(w, "There was no valid access token please request and try again.")
	}
	url := "https://api.coinbase.com/v2/accounts"
	var bearer = "Bearer " + accessToken

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", bearer)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response: ", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error Reading response bytes: ", err)
	}
	fmt.Fprint(w, string([]byte(body)))

}

const (
	GrantType     = "authorization_code"
	ClientId      = "5618ef1fe1e35a25c20c326a1daead6807ee1b3273e46e58663b0913c5fcf46c"
	ClientSecret  = "3d6d47da0b3547bb8b0a53a06af982f0f1ab6ddbff1252d0c6e55c7624ba0b0b"
	RedirectURI   = "http://localhost:8080/redirect"
	TokenEndpoint = "https://api.coinbase.com/oauth/token"
)

func requestUserData(accessToken string) (jsonResponse string) {
	url := "https://api.coinbase.com/v2/user"
	var bearer = "Bearer " + accessToken

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", bearer)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response: ", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error Reading response bytes: ", err)
	}
	jsonResponse = string([]byte(body))
	return
}

func getAccessToken(w http.ResponseWriter, r *http.Request, code string) (result string) {
	requestBody, err := json.Marshal(map[string]string{
		"grant_type":    GrantType,
		"code":          code,
		"client_id":     ClientId,
		"client_secret": ClientSecret,
		"redirect_uri":  RedirectURI,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	resp, err := http.Post(TokenEndpoint, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	oAuthAccessTokenJSON := string(body)
	var oAuthTokenData OauthToken
	json.Unmarshal([]byte(oAuthAccessTokenJSON), &oAuthTokenData)
	result = oAuthTokenData.AccessToken
	return
}
