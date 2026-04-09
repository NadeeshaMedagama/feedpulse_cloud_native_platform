package services

import "testing"

func TestAuthServiceLogin(t *testing.T) {
	auth := NewAuthService("admin@feedpulse.local", "Admin123!", "secret")

	token, err := auth.Login("admin@feedpulse.local", "Admin123!")
	if err != nil {
		t.Fatalf("expected successful login, got error: %v", err)
	}
	if token == "" {
		t.Fatal("expected token, got empty string")
	}

	if _, err := auth.Login("admin@feedpulse.local", "bad-password"); err == nil {
		t.Fatal("expected invalid credential error")
	}
}
