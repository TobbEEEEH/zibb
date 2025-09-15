package auth

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"
)

// ======================
// Token (HMAC signed)
// ======================

var hmacSecret = []byte(getEnv("AUTH_SECRET", "dev-secret"))

// GenerateToken returns a signed token for userID with TTL.
func GenerateToken(userID string, ttl time.Duration) (string, error) {
	expiry := time.Now().Add(ttl).Unix()
	nonce := make([]byte, 8)
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}
	payload := fmt.Sprintf("%s|%d|%s", userID, expiry, base64.RawStdEncoding.EncodeToString(nonce))
	sig := sign(payload)
	return base64.RawStdEncoding.EncodeToString([]byte(payload)) + "." + base64.RawStdEncoding.EncodeToString(sig), nil
}

// ParseToken verifies signature and expiry, returns userID.
func ParseToken(tok string) (string, error) {
	parts := strings.Split(tok, ".")
	if len(parts) != 2 {
		return "", errors.New("bad token format")
	}

	payloadB, err := base64.RawStdEncoding.DecodeString(parts[0])
	if err != nil {
		return "", err
	}
	sig, err := base64.RawStdEncoding.DecodeString(parts[1])
	if err != nil {
		return "", err
	}

	if !hmac.Equal(sig, sign(string(payloadB))) {
		return "", errors.New("bad signature")
	}

	// payload = userID|expiry|nonce
	fields := strings.Split(string(payloadB), "|")
	if len(fields) != 3 {
		return "", errors.New("bad payload")
	}
	userID := fields[0]
	exp, err := parseInt64(fields[1])
	if err != nil {
		return "", err
	}
	if time.Now().Unix() > exp {
		return "", errors.New("token expired")
	}
	return userID, nil
}

func sign(msg string) []byte {
	h := hmac.New(sha256.New, hmacSecret)
	h.Write([]byte(msg))
	return h.Sum(nil)
}
