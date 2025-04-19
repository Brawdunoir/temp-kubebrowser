package main

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"slices"

	v1 "github.com/brawdunoir/kubebrowser/pkg/apis/kubeconfig/v1"
	v1client "github.com/brawdunoir/kubebrowser/pkg/client/listers/kubeconfig/v1"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

type enhancedContext struct {
	logger           *zap.SugaredLogger
	oauth2Config     oauth2.Config
	oauth2Verifier   *oidc.IDTokenVerifier
	session          sessions.Session
	kubeconfigLister v1client.KubeconfigLister
}

func randString(nByte int) (string, error) {
	b := make([]byte, nByte)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// Returns a subset of initial Kubeconfigs depending on the whitelist in each Kubeconfig and the
// claims (user and groups) in the idToken
func filterKubeconfig(c *gin.Context, kubeconfigs []*v1.Kubeconfig, idToken *oidc.IDToken) ([]*v1.KubeconfigSpec, error) {
	logger := c.Request.Context().Value(loggerKey).(*zap.SugaredLogger)

	logger.Debug("Entering in filterKubeconfig")

	claims := struct {
		Email  string   `json:"email"`
		Groups []string `json:"groups"`
	}{}

	if err := idToken.Claims(&claims); err != nil {
		return nil, err
	}

	logger.Debugw("Extracted from ID token", "claims", claims)

	var filtered []*v1.KubeconfigSpec
	for _, kubeconfig := range kubeconfigs {
		logger.Debugw("Start to filter the kubeconfig", "name", kubeconfig.Name, "whitelist", kubeconfig.Spec.Whitelist)
		if kubeconfig.Spec.Whitelist == nil {
			logger.Debug("Kubeconfig has no whitelist, process next one")
			filtered = append(filtered, &kubeconfig.Spec)
			continue
		}

		if slices.Contains(kubeconfig.Spec.Whitelist.Users, claims.Email) {
			logger.Debugw("Found a match on", "kubeconfig", kubeconfig.Name, "user", claims.Email)
			filtered = append(filtered, &kubeconfig.Spec)
			continue
		} else {
			logger.Debug("Did not found a match on user, continue with groups")
		}
		for _, group := range claims.Groups {
			if slices.Contains(kubeconfig.Spec.Whitelist.Groups, group) {
				logger.Debugw("Found a match on", "kubeconfig", kubeconfig.Name, "group", group)
				filtered = append(filtered, &kubeconfig.Spec)
				break
			}
		}
	}

	return filtered, nil
}

func preprareKubeconfigs(c *gin.Context, kubeconfigs []*v1.Kubeconfig) ([]*v1.KubeconfigSpec, error) {
	ec := extractFromContext(c)
	ec.logger.Debug("Entering in prepareKubeconfigs")

	rawIDToken, refreshToken := extractTokens(ec.session)

	idToken, err := ec.oauth2Verifier.Verify(c.Request.Context(), rawIDToken)
	if err != nil {
		return nil, err
	}

	kubeconfigsSpec, err := filterKubeconfig(c, kubeconfigs, idToken)
	if err != nil {
		return nil, err
	}

	user := v1.User{Name: "oidc", Users: v1.UserSpec{
		AuthProvider: v1.AuthProviderSpec{Name: "oidc", Config: v1.AuthProviderConfig{
			ClientID:     ec.oauth2Config.ClientID,
			ClientSecret: ec.oauth2Config.ClientSecret,
			IDPIssuerURL: viper.GetString("oauth2_issuer_url"),
			IDToken:      rawIDToken,
			RefreshToken: refreshToken,
		}},
	}}

	for _, k := range kubeconfigsSpec {
		k.Whitelist = nil                                     // Remove whitelist
		k.Kubeconfig.Users = nil                              // Remove all users
		k.Kubeconfig.Users = append(k.Kubeconfig.Users, user) // Put user created before
		k.Kubeconfig.Contexts = k.Kubeconfig.Contexts[:1]     // Keep first context only
		k.Kubeconfig.Contexts[0].Context.User = user.Name     // Put same name as user
	}

	return kubeconfigsSpec, nil
}

func extractFromContext(c *gin.Context) enhancedContext {
	ec := enhancedContext{
		logger:         c.Request.Context().Value(loggerKey).(*zap.SugaredLogger),
		oauth2Config:   c.Request.Context().Value(oauth2ConfigKey).(oauth2.Config),
		oauth2Verifier: c.Request.Context().Value(oauth2VerifierKey).(*oidc.IDTokenVerifier),
		session:        sessions.Default(c),
	}
	kl := c.Request.Context().Value(kubeconfigListerKey)
	if kl == nil {
		return ec
	}
	ec.kubeconfigLister = kl.(v1client.KubeconfigLister)
	return ec
}

func extractTokens(session sessions.Session) (rawIDToken string, refreshToken string) {
	return session.Get(rawIDTokenKey).(string), session.Get(refreshTokenKey).(string)
}
