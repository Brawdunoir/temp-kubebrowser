package main

import (
	"context"
	"net/http"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

func setupOidc(ctx context.Context, clientID string, clientSecret string) (oauth2.Config, *oidc.IDTokenVerifier, error) {
	provider, err := oidc.NewProvider(ctx, "https://login.microsoftonline.com/a3594a5e-d561-4f1b-a566-9a93202ecf1d/v2.0")
	if err != nil {
		return oauth2.Config{}, nil, err
	}
	oidcConfig := &oidc.Config{
		ClientID: clientID,
	}
	verifier := provider.Verifier(oidcConfig)

	config := oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     provider.Endpoint(),
		RedirectURL:  "http://localhost:8080/auth/callback",
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email", oidc.ScopeOfflineAccess},
	}
	return config, verifier, nil
}

func refreshToken(ctx context.Context, config oauth2.Config, refreshToken string) (*oauth2.Token, error) {
	tokenSource := config.TokenSource(ctx, &oauth2.Token{RefreshToken: refreshToken})
	newToken, err := tokenSource.Token()
	if err != nil {
		return nil, err
	}
	return newToken, nil
}

func redirectToOIDCLogin(c *gin.Context, config oauth2.Config) {
	state, err := randString(16)
	if err != nil {
		c.String(http.StatusInternalServerError, "Internal error")
		c.Abort()
		return
	}
	nonce, err := randString(16)
	if err != nil {
		c.String(http.StatusInternalServerError, "Internal error")
		c.Abort()
		return
	}
	setCallbackCookie(c, "state", state)
	setCallbackCookie(c, "nonce", nonce)

	// Redirect to OIDC login
	c.Redirect(http.StatusFound, config.AuthCodeURL(state, oidc.Nonce(nonce)))
	c.Abort()
}

func handleOAuth2Callback(ctx context.Context, config oauth2.Config, verifier *oidc.IDTokenVerifier) func(*gin.Context) {
	return func(c *gin.Context) {
		state, err := c.Cookie("state")
		if err != nil {
			c.String(http.StatusBadRequest, "state not found")
			return
		}
		if c.Query("state") != state {
			c.String(http.StatusBadRequest, "state did not match")
			return
		}

		oauth2Token, err := config.Exchange(ctx, c.Query("code"))
		if err != nil {
			c.String(http.StatusInternalServerError, "Failed to exchange token: "+err.Error())
			return
		}
		rawIDToken, ok := oauth2Token.Extra("id_token").(string)
		if !ok {
			c.String(http.StatusInternalServerError, "No id_token field in oauth2 token.")
			return
		}
		idToken, err := verifier.Verify(ctx, rawIDToken)
		if err != nil {
			c.String(http.StatusInternalServerError, "Failed to verify ID Token: "+err.Error())
			return
		}

		nonce, err := c.Cookie("nonce")
		if err != nil {
			c.String(http.StatusBadRequest, "nonce not found")
			return
		}
		if idToken.Nonce != nonce {
			c.String(http.StatusBadRequest, "nonce did not match")
			return
		}

		setCallbackCookie(c, "id_token", rawIDToken)
		setCallbackCookie(c, "refresh_token", oauth2Token.RefreshToken)

		c.Redirect(http.StatusFound, "/")
	}
}
