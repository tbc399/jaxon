package plaid

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"log/slog"
	"math/big"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/plaid/plaid-go/v37/plaid"
)

// put this somewhere else
var cachedKey *plaid.JWKPublicKey

func Router() *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("POST /plaid", handlePlaidHooks)
	return router
}

func handlePlaidHooks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	plaidClient := ctx.Value("plaidClient").(*plaid.APIClient)

	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("Could not read the webhook body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	verified := verifySignature(string(body), &r.Header, plaidClient, &ctx)
	if !verified {
		slog.Warn("Signature verification failed for plaid webhook")
		w.WriteHeader(http.StatusBadRequest)
	}

	slog.Info(string(body))

	w.WriteHeader(http.StatusOK)
}

func verifySignature(webhookBody string, headers *http.Header, client *plaid.APIClient, ctx *context.Context) bool {
	// Extract the signed JWT from the webhook header
	tokenString := headers.Get("plaid-verification")

	// Parse token string, but don't validate token yet
	token, _, err := new(jwt.Parser).ParseUnverified(
		tokenString,
		jwt.MapClaims{},
	)
	if err != nil {
		return false
	}

	// Return error if alg in header is not "ES256"
	if token.Method.Alg() != "ES256" {
		return false
	}

	// Extract the key ID (kid) from the token header
	kid := token.Header["kid"].(string)

	// Fetch key if not already cached
	if cachedKey == nil {
		webhookRequest := *plaid.NewWebhookVerificationKeyGetRequest(kid)
		webhookResponse, _, respErr := client.PlaidApi.WebhookVerificationKeyGet(*ctx).WebhookVerificationKeyGetRequest(webhookRequest).Execute()

		if respErr != nil {
			fmt.Println(respErr)
			return false
		}

		// Cache the key
		key := webhookResponse.GetKey()
		cachedKey = &key
	}

	// If key is still not set, verification fails
	if cachedKey == nil {
		return false
	}

	// Signing key must be an ecdsa.PublicKey struct
	publicKey := new(ecdsa.PublicKey)
	publicKey.Curve = elliptic.P256()

	x, err := base64.URLEncoding.DecodeString(cachedKey.X + "=")
	if err != nil {
		return false
	}
	xc := new(big.Int)
	publicKey.X = xc.SetBytes(x)

	y, _ := base64.URLEncoding.DecodeString(cachedKey.Y + "=")
	if err != nil {
		return false
	}
	yc := new(big.Int)
	publicKey.Y = yc.SetBytes(y)

	// Verify the token signature
	_, err = jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return publicKey, nil
	})
	if err != nil {
		return false
	}

	// Verify claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return false
	}

	// Verify iat (issued at) claim - token should not be older than 5 minutes
	if iat, ok := claims["iat"].(float64); ok {
		issuedAt := time.Unix(int64(iat), 0)
		if time.Since(issuedAt) > 5*time.Minute {
			return false
		}
	} else {
		return false
	}

	// Verify request_body_sha256 claim matches the webhook body
	if bodyHash, ok := claims["request_body_sha256"].(string); ok {
		hasher := sha256.New()
		hasher.Write([]byte(webhookBody))
		calculatedhash := hex.EncodeToString(hasher.Sum(nil))

		if bodyHash != calculatedhash {
			return false
		}
	} else {
		return false
	}

	return true
}
