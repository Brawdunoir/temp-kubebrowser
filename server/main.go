package main

import (
	"time"

	"github.com/brawdunoir/kubebrowser/pkg/signals"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog/v2"

	clientset "github.com/brawdunoir/kubebrowser/pkg/client/clientset/versioned"

	informers "github.com/brawdunoir/kubebrowser/pkg/client/informers/externalversions"
)

// var (
// 	clientID     = os.Getenv("OAUTH2_CLIENT_ID")
// 	clientSecret = os.Getenv("OAUTH2_CLIENT_SECRET")
// )

// func randString(nByte int) (string, error) {
// 	b := make([]byte, nByte)
// 	if _, err := io.ReadFull(rand.Reader, b); err != nil {
// 		return "", err
// 	}
// 	return base64.RawURLEncoding.EncodeToString(b), nil
// }

// func setCallbackCookie(c *gin.Context, name, value string) {
// 	cookie := &http.Cookie{
// 		Name:     name,
// 		Value:    value,
// 		MaxAge:   int(time.Hour.Seconds()),
// 		Secure:   c.Request.TLS != nil,
// 		HttpOnly: true,
// 	}
// 	http.SetCookie(c.Writer, cookie)
// }

func main() {
	klog.InitFlags(nil)

	// set up signals so we handle the shutdown signal gracefully
	ctx := signals.SetupSignalHandler()
	logger := klog.FromContext(ctx)

	// creates the in-cluster config
	cfg, err := rest.InClusterConfig()
	if err != nil {
		logger.Error(err, "Error building kubeconfig")
		klog.FlushAndExit(klog.ExitFlushTimeout, 1)
	}

	exampleClient, err := clientset.NewForConfig(cfg)
	if err != nil {
		logger.Error(err, "Error building kubernetes clientset")
		klog.FlushAndExit(klog.ExitFlushTimeout, 1)
	}

	// Create the informer factory
	kubeInformerFactory := informers.NewSharedInformerFactory(exampleClient, time.Second*30)

	// Get the lister for Kubeconfigs
	kubeconfigLister := kubeInformerFactory.Kubeconfig().V1().Kubeconfigs().Lister()

	// Start the informer factory
	stopCh := make(chan struct{})
	kubeInformerFactory.Start(stopCh)

	// Wait for the caches to sync
	if !cache.WaitForCacheSync(stopCh, kubeInformerFactory.Kubeconfig().V1().Kubeconfigs().Informer().HasSynced) {
		logger.Error(nil, "Failed to sync caches")
		klog.FlushAndExit(klog.ExitFlushTimeout, 1)
	}

	// Use the lister to get the Kubeconfig resource
	kubeconfig, err := kubeconfigLister.Kubeconfigs("default").Get("example-kubeconfig")
	if err != nil {
		logger.Error(err, "Failed to get example-kubeconfig")
		klog.FlushAndExit(klog.ExitFlushTimeout, 1)
	}

	logger.Info("I fetched example-kubeconfig once:", "example-kubeconfig", kubeconfig)
	// time.Sleep(time.Hour * 1)

	for {
		kubeconfigs, err := kubeconfigLister.Kubeconfigs("default").List(labels.NewSelector())
		if err != nil {
			logger.Error(err, "Error listing kubeconfigs")
			klog.FlushAndExit(klog.ExitFlushTimeout, 1)
		}
		for n, kubeconfig := range kubeconfigs {

			logger.Info("this is kubeconfig:", "number", n, "kubeconfig", kubeconfig)
		}
		logger.Info("I have listed all kubeconfigs")
		time.Sleep(time.Second * 10)
	}
}

// ctx := context.Background()

// provider, err := oidc.NewProvider(ctx, "https://login.microsoftonline.com/a3594a5e-d561-4f1b-a566-9a93202ecf1d/v2.0")
// if err != nil {
// 	log.Fatal(err)
// }
// oidcConfig := &oidc.Config{
// 	ClientID: clientID,
// }
// verifier := provider.Verifier(oidcConfig)

// config := oauth2.Config{
// 	ClientID:     clientID,
// 	ClientSecret: clientSecret,
// 	Endpoint:     provider.Endpoint(),
// 	RedirectURL:  "http://localhost:8080/auth/callback",
// 	Scopes:       []string{oidc.ScopeOpenID, "profile", "email", oidc.ScopeOfflineAccess},
// }

// r := gin.Default()

// r.GET("/", func(c *gin.Context) {
// 	state, err := randString(16)
// 	if err != nil {
// 		c.String(http.StatusInternalServerError, "Internal error")
// 		return
// 	}
// 	nonce, err := randString(16)
// 	if err != nil {
// 		c.String(http.StatusInternalServerError, "Internal error")
// 		return
// 	}
// 	setCallbackCookie(c, "state", state)
// 	setCallbackCookie(c, "nonce", nonce)

// 	c.Redirect(http.StatusFound, config.AuthCodeURL(state, oidc.Nonce(nonce)))
// })

// r.GET("/auth/callback", func(c *gin.Context) {
// 	state, err := c.Cookie("state")
// 	if err != nil {
// 		c.String(http.StatusBadRequest, "state not found")
// 		return
// 	}
// 	if c.Query("state") != state {
// 		c.String(http.StatusBadRequest, "state did not match")
// 		return
// 	}

// 	oauth2Token, err := config.Exchange(ctx, c.Query("code"))
// 	if err != nil {
// 		c.String(http.StatusInternalServerError, "Failed to exchange token: "+err.Error())
// 		return
// 	}
// 	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
// 	if !ok {
// 		c.String(http.StatusInternalServerError, "No id_token field in oauth2 token.")
// 		return
// 	}
// 	idToken, err := verifier.Verify(ctx, rawIDToken)
// 	if err != nil {
// 		c.String(http.StatusInternalServerError, "Failed to verify ID Token: "+err.Error())
// 		return
// 	}

// 	nonce, err := c.Cookie("nonce")
// 	if err != nil {
// 		c.String(http.StatusBadRequest, "nonce not found")
// 		return
// 	}
// 	if idToken.Nonce != nonce {
// 		c.String(http.StatusBadRequest, "nonce did not match")
// 		return
// 	}

// 	oauth2Token.AccessToken = "*REDACTED*"

// 	resp := struct {
// 		RefreshToken string
// 		IDToken      string
// 	}{oauth2Token.RefreshToken, rawIDToken}

// 	data, err := json.Marshal(resp)
// 	if err != nil {
// 		c.String(http.StatusInternalServerError, err.Error())
// 		return
// 	}
// 	c.Data(http.StatusOK, "application/json", data)
// })

// log.Fatal(r.Run())
