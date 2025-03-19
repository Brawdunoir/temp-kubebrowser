package main

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"slices"

	v1 "github.com/brawdunoir/kubebrowser/pkg/apis/kubeconfig/v1"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func randString(nByte int) (string, error) {
	b := make([]byte, nByte)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// Returns a subset of initial Kubeconfigs depending on the whitelist in each Kubeconfig and the
// claims (user and groups) in the idToken
func filterKubeconfig(c *gin.Context, kubeconfigs []*v1.Kubeconfig, idToken *oidc.IDToken) ([]*v1.KubeconfigData, error) {
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

	var filtered []*v1.KubeconfigData
	for _, kubeconfig := range kubeconfigs {
		logger.Debugw("Start to filter the kubeconfig", "name", kubeconfig.Name, "whitelist", kubeconfig.Spec.Whitelist)

		if slices.Contains(kubeconfig.Spec.Whitelist.Users, claims.Email) {
			logger.Debugw("Found a match on", "kubeconfig", kubeconfig.Name, "user", claims.Email)
			filtered = append(filtered, &kubeconfig.Spec.Kubeconfig)
			continue
		} else {
			logger.Debug("Did not found a match on user, continue with groups")
		}
		for _, group := range claims.Groups {
			if slices.Contains(kubeconfig.Spec.Whitelist.Groups, group) {
				logger.Debugw("Found a match on", "kubeconfig", kubeconfig.Name, "group", group)
				filtered = append(filtered, &kubeconfig.Spec.Kubeconfig)
				break
			}
		}
	}

	return filtered, nil
}
