// Package auth serves middleware, sessions and authentication
package auth

import (
	"fmt"
	"os"
)

// ======================
// Helpers
// ======================

func getEnv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func parseInt64(s string) (int64, error) {
	var i int64
	_, err := fmt.Sscan(s, &i)
	return i, err
}
