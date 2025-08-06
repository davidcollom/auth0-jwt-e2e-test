package main

import (
	"context"
	"fmt"

	"golang.org/x/oauth2/clientcredentials"
)

func GetTokenUsingOAuth2(ctx context.Context, domain, clientID, clientSecret, audience string) (string, error) {
	conf := &clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     fmt.Sprintf("https://%s/oauth/token", domain),
		EndpointParams: map[string][]string{
			"audience": {audience}, // 🔥 required by Auth0
		},
	}

	token, err := conf.Token(ctx)
	if err != nil {
		return "", err
	}

	return token.AccessToken, nil
}
