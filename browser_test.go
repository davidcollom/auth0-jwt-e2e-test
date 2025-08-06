package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"testing"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/stretchr/testify/require"
)

const BaseURL = "https://httpbin.org"

func TestHomepage(t *testing.T) {
	// Skip if Auth0 tests are disabled or JWT token is not available
	if os.Getenv("SKIP_AUTH0_TESTS") == "true" || JWTToken == "" {
		t.Skip("Skipping browser test - Auth0 integration not available")
	}

	t.Logf("Starting homepage test...")
	ctx, cancel := chromedp.NewContext(t.Context())
	defer cancel()

	t.Logf("Beginning capture setup")

	done := make(chan bool)
	var requestID network.RequestID

	chromedp.ListenTarget(ctx, func(v interface{}) {
		switch ev := v.(type) {
		case *network.EventRequestWillBeSent:
			log.Printf("EventRequestWillBeSent: %v: %v", ev.RequestID, ev.Request.URL)
			requestID = ev.RequestID
		case *network.EventLoadingFinished:
			log.Printf("EventLoadingFinished: %v", ev.RequestID)
			if ev.RequestID == requestID {
				close(done)
			}
		}
	})

	err := chromedp.Run(ctx,
		network.Enable(),
		network.SetExtraHTTPHeaders(network.Headers{
			"Authorization": "Bearer " + JWTToken,
		}),
		chromedp.Navigate(BaseURL+"/bearer"),
	)
	require.NoError(t, err, "failed to navigate to /bearer")

	<-done // wait for the response body to be retrieved

	var buf []byte
	if err := chromedp.Run(ctx, chromedp.ActionFunc(func(ctx context.Context) error {
		var err error
		buf, err = network.GetResponseBody(requestID).Do(ctx)
		return err
	})); err != nil {
		log.Fatal(err)
	}
	require.NotEmpty(t, buf, "response body should not be empty")

	t.Logf("✅ JSON body from /bearer: %s", buf)
	var jsonBody JsonResponse
	require.NoError(t, json.Unmarshal(buf, &jsonBody))

	require.True(t, jsonBody.Authenticated, "authenticated key missing or wrong type")
	require.NotEmpty(t, jsonBody.Token, "token key missing or wrong type")
}

type JsonResponse struct {
	Authenticated bool   `json:"authenticated,omitempty"`
	Token         string `json:"token,omitempty"`
}
