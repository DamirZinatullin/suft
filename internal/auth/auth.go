package auth

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

const (
	AuthURL    = "security/authenticate"
	BaseURL    = "https://dev.gnivc.ru/tools/suft/api/v1/"
	RefreshURL = "security/refresh-token"
)

type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type Options struct {
	SuftAPIURL  string
	HttpTimeout time.Duration
}

func Authenticate(email string, password string, options *Options) (*Token, error) {
	baseURL := BaseURL
	httpTimeout := 2 * time.Second
	if options != nil {
		if options.SuftAPIURL != "" {
			baseURL = options.SuftAPIURL
		}
		if options.HttpTimeout != 0 {
			httpTimeout = options.HttpTimeout
		}
	}

	cli := &http.Client{
		Timeout: httpTimeout,
	}
	token := &Token{}

	reqBody := bytes.NewBuffer([]byte(fmt.Sprintf(`{"username":"%s","password":"%s"}`, email, password)))
	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprint(baseURL, AuthURL),
		reqBody,
	)
	if err != nil {
		log.Println("Auth: unable to create new request:", err)
		return token, err
	}
	req.Header.Add("Auth-method", "Password")
	req.Header.Add("Content-Type", "application/json; charset=UTF-8")

	resp, err := cli.Do(req)
	if err != nil {
		log.Println("Auth: unable to get http response:", err)
		return token, err
	}

	defer resp.Body.Close()

	respB, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Auth: unable to read authentification tokens:", err)
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(string(respB))
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

func Refresh(refreshToken string, options *Options) (*Token, error) {
	baseURL := BaseURL
	httpTimeout := 2 * time.Second
	if options != nil {
		if options.SuftAPIURL != "" {
			baseURL = options.SuftAPIURL
		}
		if options.HttpTimeout != 0 {
			httpTimeout = options.HttpTimeout
		}
	}

	cli := &http.Client{
		Timeout: httpTimeout,
	}
	token := &Token{}

	req, err := http.NewRequest(
		http.MethodPut,
		fmt.Sprint(baseURL, RefreshURL),
		nil,
	)
	if err != nil {
		log.Println("Refresh: unable to create new request:", err)
		return nil, err
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
		log.Println("Refresh: unable to get http response:", err)
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("unable to refresh tokens")
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
