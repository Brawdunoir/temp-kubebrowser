package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"time"

	v1 "github.com/brawdunoir/kubebrowser/pkg/client/listers/kubeconfig/v1"
	"github.com/brawdunoir/kubebrowser/pkg/signals"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	"k8s.io/apimachinery/pkg/labels"
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

	// Create controller lister for Kubeconfigs CRD
	kubeconfigLister, err := setupKubeconfigLister(ctx)
	if err != nil {
		logger.Error(err, "Cannot setup kubeconfig lister")
		klog.FlushAndExit(klog.ExitFlushTimeout, 1)
	}

	// Create OIDC related config and verifier
	config, verifier, err := setupOidc(ctx, clientID, clientSecret)
	if err != nil {
		logger.Error(err, "Failed to setup Oidc")
		klog.FlushAndExit(klog.ExitFlushTimeout, 1)
	}

	// Create session store
	store := memstore.NewStore([]byte(sessionSecret))

	router := gin.Default()
	router.Use(sessions.Sessions("kubebrowser_session", store))
	router.Use(AuthMiddleware(verifier, config))

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})
	router.GET(callbackRoute, handleOAuth2Callback(config, verifier))
	router.GET("/api/kubeconfigs", handleGetKubeconfigs(verifier, kubeconfigLister))

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

func handleGetKubeconfigs(verifier *oidc.IDTokenVerifier, kl v1.KubeconfigLister) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := c.Request.Context().Value(loggerKey).(klog.Logger)
		session := sessions.Default(c)

		logger.Info("Getting kubeconfigs")

		kubeconfigs, err := kl.Kubeconfigs("default").List(labels.Everything())
		if err != nil {
			logger.Error(err, "Error listing kubeconfigs")
			klog.FlushAndExit(klog.ExitFlushTimeout, 1)
		}

		rawIDToken := session.Get(rawIDTokenKey)
		idToken, err := verifier.Verify(c.Request.Context(), rawIDToken.(string))

		filteredKubeconfigs := filterKubeconfig(kubeconfigs, idToken)

		for _, k := range filteredKubeconfigs {
			logger.Info("This is a kubeconfig you are allowed to see", k.Name, k.Spec.Whitelist)
		}
	}
}
