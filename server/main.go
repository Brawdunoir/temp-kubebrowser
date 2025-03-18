package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/brawdunoir/kubebrowser/pkg/signals"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	"k8s.io/klog/v2"
)

type contextKey string

const (
	loggerKey     contextKey = "logger"
	callbackRoute string     = "/auth/callback"
	sessionSecret string     = "secret"
)

func main() {
	klog.InitFlags(nil)

	// set up signals so we handle the shutdown signal gracefully
	ctx := signals.SetupSignalHandler()
	logger := klog.FromContext(ctx)

	// Add logger to context
	ctx = context.WithValue(ctx, loggerKey, logger)

	config, verifier, err := setupOidc(ctx, clientID, clientSecret)
	if err != nil {
		logger.Error(err, "Failed to setup Oidc")
		klog.FlushAndExit(klog.ExitFlushTimeout, 1)
	}
	store := memstore.NewStore([]byte(sessionSecret))

	router := gin.Default()

	router.Use(sessions.Sessions("kubebrowser_session", store))
	router.Use(AuthMiddleware(verifier, config))

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})
	router.GET(callbackRoute, handleOAuth2Callback(config, verifier))
	router.GET("/api/kubeconfigs")

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
		BaseContext: func(net.Listener) context.Context {
			return ctx
		},
	}

	// Run the server in a goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for the shutdown signal
	<-ctx.Done()

	// Gracefully shut down the server
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error(err, "Server forced to shutdown")
	}

	logger.Info("Server exiting")
}
