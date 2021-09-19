package auth

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

const authURL = "security/authenticate"
const baseURL = "https://dev.gnivc.ru/tools/suft/api/v1/"
const refreshURL = "security/refresh-token"

type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func Authenticate(email string, password string) (*Token, error) {
	cli := &http.Client{
		Timeout: 10 * time.Second,
	}
	token := &Token{}

	reqBody := bytes.NewBuffer([]byte(fmt.Sprintf(`{"username":"%s","password":"%s"}`, email, password)))
	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprint(baseURL, authURL),
		reqBody,
	)
	if err != nil {
		log.Println(err)
		return token, err
	}
	req.Header.Add("Auth-method", "Password")
	req.Header.Add("Content-Type", "application/json; charset=UTF-8")

	resp, err := cli.Do(req)
	if err != nil {
		return token, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return token, errors.New("unable to get authentification tokens")
	}

	for _, cookie := range resp.Cookies() {
		if cookie.Name == "Access-token" {
			token.AccessToken = cookie.Value
		}
		if cookie.Name == "Refresh-token" {
			token.RefreshToken = cookie.Value
		}
	}

	return token, nil
}

func Refresh(refreshToken string) (*Token, error) {
	cli := &http.Client{
		Timeout: 10 * time.Second,
	}
	token := &Token{}

	req, err := http.NewRequest(
		http.MethodPut,
		fmt.Sprint(baseURL, refreshURL),
		nil,
	)
	if err != nil {
		log.Println(err)
		return token, err
	}
	req.Header.Add("Auth-method", "Password")
	req.Header.Add("Content-Type", "application/json; charset=UTF-8")
	cookie1 := &http.Cookie{
		Name:  "Refresh-token",
		Value: refreshToken,
	}
	req.AddCookie(cookie1)

	resp, err := cli.Do(req)
	if err != nil {
		return token, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return token, errors.New("unable to get authentification tokens")
	}

	for _, cookie := range resp.Cookies() {
		if cookie.Name == "Access-token" {
			token.AccessToken = cookie.Value
		}
		if cookie.Name == "Refresh-token" {
			token.RefreshToken = cookie.Value
		}
	}

	return token, nil
}
