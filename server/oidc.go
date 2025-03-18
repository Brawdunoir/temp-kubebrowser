package main

import (
	"context"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"k8s.io/klog/v2"
)

var (
	clientID        = os.Getenv("OAUTH2_CLIENT_ID")
	clientSecret    = os.Getenv("OAUTH2_CLIENT_SECRET")
	rawIDTokenKey   = "id_token"
	refreshTokenKey = "refresh_token"
)

func setCallbackCookie(c *gin.Context, name, value string) {
	logger := c.Request.Context().Value(loggerKey).(klog.Logger)
	c.SetCookie(
		name,
		value,
		int(time.Hour.Seconds()),              // MaxAge in seconds
		"/",                                   // Path
		strings.Split(c.Request.Host, ":")[0], // Domain (empty for default)
		c.Request.TLS != nil,                  // Secure (set to false for local development)
		true,                                  // HttpOnly
	)
	logger.Info("Cookie is set", "cookie", name)
}

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
		RedirectURL:  "http://localhost:8080" + callbackRoute,
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email", oidc.ScopeOfflineAccess},
	}
	return config, verifier, nil
}

func refreshTokens(ctx context.Context, config oauth2.Config, refreshToken string) (*oauth2.Token, error) {
	tokenSource := config.TokenSource(ctx, &oauth2.Token{RefreshToken: refreshToken})
	newToken, err := tokenSource.Token()
	if err != nil {
		return nil, err
	}
	return newToken, nil
}

func redirectToOIDCLogin(c *gin.Context, config oauth2.Config) {

	logger := c.Request.Context().Value(loggerKey).(klog.Logger)

	logger.Info("Entering in redirectToOIDCLogin")

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
}

func handleOAuth2Callback(config oauth2.Config, verifier *oidc.IDTokenVerifier) func(*gin.Context) {
	return func(c *gin.Context) {
		logger := c.Request.Context().Value(loggerKey).(klog.Logger)
		session := sessions.Default(c)

		logger.Info("Entering in handleOAuth2Callback")

		// Retrieve and validate state
		state, err := c.Cookie("state")
		if err != nil {
			logger.Error(err, "State cookie not found")
			c.String(http.StatusBadRequest, "State not found")
			return
		}
		if c.Query("state") != state {
			logger.Error(nil, "State mismatch", "expected", state, "got", c.Query("state"))
			c.String(http.StatusBadRequest, "State did not match")
			return
		}

		// Exchange code for token
		oauth2Token, err := config.Exchange(c.Request.Context(), c.Query("code"))
		if err != nil {
			logger.Error(err, "Failed to exchange token")
			c.String(http.StatusInternalServerError, "Failed to exchange token: "+err.Error())
			return
		}

		// Retrieve and validate nonce
		rawIDToken, ok := oauth2Token.Extra("id_token").(string)
		if !ok {
			logger.Error(nil, "No id_token field in oauth2 token")
			c.String(http.StatusInternalServerError, "No id_token field in oauth2 token")
			return
		}
		idToken, err := verifier.Verify(c.Request.Context(), rawIDToken)
		if err != nil {
			logger.Error(err, "Failed to verify ID Token")
			c.String(http.StatusInternalServerError, "Failed to verify ID Token: "+err.Error())
			return
		}

		nonce, err := c.Cookie("nonce")
		if err != nil {
			logger.Error(err, "Nonce cookie not found")
			c.String(http.StatusBadRequest, "Nonce not found")
			return
		}
		if idToken.Nonce != nonce {
			logger.Error(nil, "Nonce mismatch", "expected", nonce, "got", idToken.Nonce)
			c.String(http.StatusBadRequest, "Nonce did not match")
			return
		}

		session.Set(rawIDTokenKey, rawIDToken)
		session.Set(refreshTokenKey, oauth2Token.RefreshToken)
		err = session.Save()
		if err != nil {
			logger.Error(err, "Cannot save session")
			c.String(http.StatusInternalServerError, "Cannot save session")
		}
		c.Redirect(http.StatusFound, "/")
	}
}

func AuthMiddleware(verifier *oidc.IDTokenVerifier, config oauth2.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := c.Request.Context().Value(loggerKey).(klog.Logger)
		session := sessions.Default(c)

		logger.Info("Entering in AuthMiddleware")
		// Skip authentication for callback route
		if c.Request.URL.Path == callbackRoute {
			logger.Info("Skip auth because user hitting callback route")
			c.Next()
		}

		// Retrieve ID token from session
		rawIDToken := session.Get(rawIDTokenKey)
		if rawIDToken == nil {
			logger.Info("ID token missing, redirecting to login")
			redirectToOIDCLogin(c, config)
			return
		}

		// Verify ID token
		idToken, err := verifier.Verify(c.Request.Context(), rawIDToken.(string))
		if err != nil || idToken.Expiry.Before(time.Now()) {
			logger.Info("ID token expired or invalid, attempting to refresh")

			// Retrieve refresh token
			refreshToken := session.Get(refreshTokenKey)
			if refreshToken == nil {
				logger.Info("Refresh token cookie missing, redirecting to login")
				redirectToOIDCLogin(c, config)
				return
			}

			// Refresh tokens
			newToken, err := refreshTokens(c.Request.Context(), config, refreshToken.(string))
			if err != nil {
				logger.Error(err, "Failed to refresh token, redirecting to login")
				redirectToOIDCLogin(c, config)
				return
			}

			// Verify new ID token
			_, err = verifier.Verify(c.Request.Context(), newToken.Extra("id_token").(string))
			if err != nil {
				logger.Error(err, "Failed to verify refreshed ID token, redirecting to login")
				redirectToOIDCLogin(c, config)
				return
			}

			// Update session with new tokens
			session.Set(rawIDTokenKey, newToken.Extra("id_token").(string))
			session.Set(refreshTokenKey, newToken.RefreshToken)
			err = session.Save()
			if err != nil {
				logger.Error(err, "Cannot save session")
				c.String(http.StatusInternalServerError, "Cannot save session")
			}
		}

		// Token is valid, proceed with the request
		logger.Info("Token is valid, proceed with the request")
		c.Next()
	}
}
