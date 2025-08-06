package main

import (
	"context"
	"log"
	"os"
	"testing"

	"gopkg.in/auth0.v1/management"
)

var (
	JWTToken string
	ClientID string
)

func TestMain(m *testing.M) {
	domain := os.Getenv("AUTH0_DOMAIN") // e.g., "dev-xxxxx.auth0.com"
	ClientID = os.Getenv("AUTH0_MANAGEMENT_CLIENT_ID")
	clientSecret := os.Getenv("AUTH0_MANAGEMENT_CLIENT_SECRET")

	// Create a new Auth0 management client
	mgmt, err := management.New(domain, ClientID, clientSecret)
	if err != nil {
		log.Fatalf("failed to create management client: %v", err)
	}
	log.Println("Auth0 Management Client created successfully")

	log.Println("Checking for existing clients...")
	clients, err := mgmt.Client.List()
	if err != nil {
		log.Fatalf("failed to list clients: %v", err)
	}
	log.Println("Existing clients:", len(clients))
	for _, client := range clients {
		if *client.Name == "My E2E Auth App" {
			mgmt.Client.Delete(*client.ClientID)
			log.Printf("Deleted existing client: %s", *client.ClientID)
			continue
		}
	}

	// 🧪 Create a new machine-to-machine app
	app := &management.Client{
		Name:       auth0String("My E2E Auth App"),
		AppType:    auth0String("non_interactive"),
		GrantTypes: []interface{}{"client_credentials"},
	}

	err = mgmt.Client.Create(app)
	if err != nil {
		log.Printf("failed to create client: %v", err)
		return
	}
	log.Printf("🧪 Created app: ID=%s, secret=%s\n", *app.ClientID, *app.ClientSecret)

	// Create a client grant for the app
	log.Println("Creating client grant for the app...")
	err = mgmt.ClientGrant.Create(&management.ClientGrant{
		ClientID: app.ClientID,
		Audience: auth0String(BaseURL),
		Scope:    []interface{}{"read:all"}, // or whatever scopes your API needs
	})
	if err != nil {
		log.Printf("⚠️ failed to create client grant: %v", err)
	}

	// 🧹 Teardown: delete client
	defer func() {
		err := mgmt.Client.Delete(*app.ClientID)
		if err != nil {
			log.Printf("⚠️ Failed to delete app: %v", err)
		} else {
			log.Println("✅ Deleted test client")
		}
	}()

	JWTToken, err = GetTokenUsingOAuth2(
		context.Background(),
		domain,
		auth0StringValue(app.ClientID),
		auth0StringValue(app.ClientSecret),
		BaseURL,
	)
	log.Printf("JWT Token: %s", JWTToken)
	if err != nil {
		log.Printf("failed to get OAuth2 token: %v", err)
		return
	}
	log.Println("OAuth2 token retrieved successfully")

	log.Println("Auth0 App created. Starting tests...")

	m.Run()
	log.Println("Tests completed. Exiting...")

}

// Helper functions to work with Auth0's string pointers
func auth0String(v string) *string {
	return &v
}

func auth0StringValue(v *string) string {
	if v != nil {
		return *v
	}
	return ""
}
