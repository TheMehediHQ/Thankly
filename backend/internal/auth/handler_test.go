package auth

import (
	"testing"
)

func TestGenerateTokens(t *testing.T) {
	userID := "test-user-123"

	token, refreshToken, err := generateTokens(userID)
	if err != nil {
		t.Fatalf("generateTokens() error = %v", err)
	}

	if token == "" {
		t.Error("generateTokens() token is empty")
	}

	if refreshToken == "" {
		t.Error("generateTokens() refreshToken is empty")
	}

	if token == refreshToken {
		t.Error("generateTokens() token and refreshToken should be different")
	}
}
