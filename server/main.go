package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/brawdunoir/kubebrowser/pkg/signals"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"k8s.io/klog/v2"
)

var (
	clientID     = os.Getenv("OAUTH2_CLIENT_ID")
	clientSecret = os.Getenv("OAUTH2_CLIENT_SECRET")
)

func main() {
	klog.InitFlags(nil)

	// set up signals so we handle the shutdown signal gracefully
	ctx := signals.SetupSignalHandler()
	logger := klog.FromContext(ctx)

	// kubeconfigLister, err := setupKubeconfigLister(ctx)
	// if err != nil {
	// 	logger.Error(err, "Error creating kubeconfigLister")
	// 	klog.FlushAndExit(klog.ExitFlushTimeout, 1)
	// }

	config, verifier, err := setupOidc(ctx, clientID, clientSecret)
	if err != nil {
		logger.Error(err, "Failed to setup Oidc")
		klog.FlushAndExit(klog.ExitFlushTimeout, 1)
	}

	router := gin.Default()

	router.Use(AuthMiddleware(verifier, config))

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Welcome to the authenticated app!")
	})

	router.GET("/auth/callback", handleOAuth2Callback(ctx, config, verifier))

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
		BaseContext: func(net.Listener) context.Context {
			return ctx
		},
	}
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}
}

func AuthMiddleware(verifier *oidc.IDTokenVerifier, config oauth2.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		idTokenCookie, err := c.Cookie("id_token")
		if err != nil || idTokenCookie == "" {
			redirectToOIDCLogin(c, config)
			return
		}

		idToken, err := verifier.Verify(c.Request.Context(), idTokenCookie)
		if err != nil || idToken.Expiry.Before(time.Now()) {
			// Attempt to refresh the token
			refreshTokenCookie, err := c.Cookie("refresh_token")
			if err != nil || refreshTokenCookie == "" {
				redirectToOIDCLogin(c, config)
				return
			}

			newToken, err := refreshToken(c.Request.Context(), config, refreshTokenCookie)
			if err != nil {
				redirectToOIDCLogin(c, config)
				return
			}

			// Update cookies with the new tokens
			setCallbackCookie(c, "id_token", newToken.Extra("id_token").(string))
			setCallbackCookie(c, "refresh_token", newToken.RefreshToken)

			// Verify the new ID token
			_, err = verifier.Verify(c.Request.Context(), newToken.Extra("id_token").(string))
			if err != nil {
				redirectToOIDCLogin(c, config)
				return
			}
		}

		// Token is valid, proceed with the request
		c.Next()
	}
}
