package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestMakeJWTAndValidateJWT(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "supersecretkey"
	expiresIn := time.Minute * 5

	tokenString, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT() error = %v, wantErr %v", err, false)
	}
	if tokenString == "" {
		t.Errorf("MakeJWT() returned empty token string")
	}

	parsedUserID, err := ValidateJWT(tokenString, tokenSecret)
	if err != nil {
		t.Fatalf("ValidateJWT() error = %v, wantErr %v", err, false)
	}
	if parsedUserID != userID {
		t.Errorf("ValidateJWT() = %v, want %v", parsedUserID, userID)
	}
}

func TestMakeJWT_EmptySecret(t *testing.T) {
	userID := uuid.New()
	tokenSecret := ""
	expiresIn := time.Minute * 5

	_, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err == nil {
		t.Errorf("MakeJWT() expected error for empty secret, got nil")
	}
}

func TestMakeJWT_NilUserID(t *testing.T) {
	userID := uuid.Nil
	tokenSecret := "supersecretkey"
	expiresIn := time.Minute * 5
	_, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err == nil {
		t.Errorf("MakeJWT() expected error for nil user ID, got nil")
	}
	// Check if the error message is as expected
	if err.Error() != "user ID cannot be nil" {
		t.Errorf("MakeJWT() expected error 'user ID cannot be nil', got %v", err)
	}
}

func TestValidateJWT_InvalidToken(t *testing.T) {
	tokenSecret := "supersecretkey"
	invalidToken := "this.is.not.a.valid.jwt"

	_, err := ValidateJWT(invalidToken, tokenSecret)
	if err == nil {
		t.Errorf("ValidateJWT() expected error for invalid token, got nil")
	}
}

func TestValidateJWT_WrongSecret(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "supersecretkey"
	wrongSecret := "wrongsecret"
	expiresIn := time.Minute * 5

	tokenString, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT() error = %v, wantErr %v", err, false)
	}

	_, err = ValidateJWT(tokenString, wrongSecret)
	if err == nil {
		t.Errorf("ValidateJWT() expected error for wrong secret, got nil")
	}
}

func TestValidateJWT_ExpiredToken(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "supersecretkey"
	expiresIn := -time.Minute // already expired

	tokenString, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT() error = %v, wantErr %v", err, false)
	}

	_, err = ValidateJWT(tokenString, tokenSecret)
	if err == nil {
		t.Errorf("ValidateJWT() expected error for expired token, got nil")
	}
}

func TestValidateJWT_EmptyToken(t *testing.T) {
	tokenString := ""
	tokenSecret := "supersecretkey"
	_, err := ValidateJWT(tokenString, tokenSecret)
	if err == nil {
		t.Errorf("ValidateJWT() expected error for empty token, got nil")
	}
	// Check if the error message is as expected
	if err.Error() != "token string cannot be empty" {
		t.Errorf("ValidateJWT() expected error 'token string cannot be empty', got %v", err)
	}
}
func TestValidateJWT_EmptySecret(t *testing.T) {
	tokenString := "valid"
	tokenSecret := ""
	_, err := ValidateJWT(tokenString, tokenSecret)
	if err == nil {
		t.Errorf("ValidateJWT() expected error for empty secret, got nil")
	}
	// Check if the error message is as expected
	if err.Error() != "token secret cannot be empty" {
		t.Errorf("ValidateJWT() expected error 'token secret cannot be empty', got %v", err)
	}
}

// test that two userIds do not produce the same token
func TestMakeJWT_DifferentUserIDs(t *testing.T) {
	userID1 := uuid.New()
	userID2 := uuid.New()
	tokenSecret := "supersecretkey"
	expiresIn := time.Minute * 5
	tokenString1, err := MakeJWT(userID1, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT() error = %v, wantErr %v", err, false)
	}
	tokenString2, err := MakeJWT(userID2, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT() error = %v, wantErr %v", err, false)
	}
	if tokenString1 == tokenString2 {
		t.Errorf("MakeJWT() returned same token string for different user IDs")
	}
	// Validate that the tokens are valid for their respective user IDs
	parsedUserID1, err := ValidateJWT(tokenString1, tokenSecret)
	if err != nil {
		t.Fatalf("ValidateJWT() error = %v, wantErr %v", err, false)
	}
	if parsedUserID1 != userID1 {
		t.Errorf("ValidateJWT() = %v, want %v", parsedUserID1, userID1)
	}
	parsedUserID2, err := ValidateJWT(tokenString2, tokenSecret)
	if err != nil {
		t.Fatalf("ValidateJWT() error = %v, wantErr %v", err, false)
	}
	if parsedUserID2 != userID2 {
		t.Errorf("ValidateJWT() = %v, want %v", parsedUserID2, userID2)
	}
}
