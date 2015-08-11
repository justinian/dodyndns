package main

import "golang.org/x/oauth2"

type TokenString struct {
	token string
}

func NewTokenString(token string) *TokenString {
	return &TokenString{token}
}

func (ts *TokenString) Token() (*oauth2.Token, error) {
	t := &oauth2.Token{
		AccessToken: ts.token,
	}
	return t, nil
}
