package main

import (
	"crypto/rand"
	"encoding/base64"
	"io"

	v1 "github.com/brawdunoir/kubebrowser/pkg/apis/kubeconfig/v1"
	"github.com/coreos/go-oidc/v3/oidc"
)

func randString(nByte int) (string, error) {
	b := make([]byte, nByte)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func filterKubeconfig(kubeconfigs []*v1.Kubeconfig, idToken *oidc.IDToken) []*v1.Kubeconfig {
	return kubeconfigs
}
