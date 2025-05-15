package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestValidateChirpHandler_MethodNotAllowed(t *testing.T) {
	cfg := &apiConfig{}
	req := httptest.NewRequest(http.MethodGet, "/api/validate_chirp", nil)
	rr := httptest.NewRecorder()

	cfg.validateChirpHandler(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status 405, got %d", rr.Code)
	}
	var resp map[string]string
	json.Unmarshal(rr.Body.Bytes(), &resp)
	if resp["error"] != "Method not allowed" {
		t.Errorf("expected error message, got %v", resp["error"])
	}
}

func TestValidateChirpHandler_BadJSON(t *testing.T) {
	cfg := &apiConfig{}
	body := bytes.NewBufferString("{bad json}")
	req := httptest.NewRequest(http.MethodPost, "/api/validate_chirp", body)
	rr := httptest.NewRecorder()

	cfg.validateChirpHandler(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", rr.Code)
	}
	var resp map[string]string
	json.Unmarshal(rr.Body.Bytes(), &resp)
	if resp["error"] != "Something went wrong" {
		t.Errorf("expected error message, got %v", resp["error"])
	}
}

func TestValidateChirpHandler_TooLong(t *testing.T) {
	cfg := &apiConfig{}
	longBody := make([]byte, 141)
	for i := range longBody {
		longBody[i] = 'a'
	}
	payload := map[string]string{"body": string(longBody)}
	b, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/validate_chirp", bytes.NewReader(b))
	rr := httptest.NewRecorder()

	cfg.validateChirpHandler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rr.Code)
	}
	var resp map[string]string
	json.Unmarshal(rr.Body.Bytes(), &resp)
	if resp["error"] != "Chirp is too long" {
		t.Errorf("expected error message, got %v", resp["error"])
	}
}

func TestValidateChirpHandler_ProfanityMultiple(t *testing.T) {
	cfg := &apiConfig{}
	tests := []struct {
		input    string
		expected string
	}{
		{"kerfuffle sharbert fornax", "**** **** ****"},
		{"hello kerfuffle sharbert world", "hello **** **** world"},
		{"KERFUFFLE SHARBERT FORNAX", "**** **** ****"},
		{"Kerfuffle and sharbert are here", "**** and **** are here"},
		{"no bad words here", "no bad words here"},
		{"fornax, kerfuffle!", "fornax, kerfuffle!"},
		{"sharbert.", "sharbert."},
	}

	for _, tt := range tests {
		payload := map[string]string{"body": tt.input}
		b, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPost, "/api/validate_chirp", bytes.NewReader(b))
		rr := httptest.NewRecorder()

		cfg.validateChirpHandler(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("input %q: expected status 200, got %d", tt.input, rr.Code)
		}
		var resp map[string]string
		json.Unmarshal(rr.Body.Bytes(), &resp)
		if resp["cleaned_body"] != tt.expected {
			t.Errorf("input %q: expected cleaned_body %q, got %q", tt.input, tt.expected, resp["cleaned_body"])
		}
	}
}

func TestValidateChirpHandler_ValidChirp(t *testing.T) {
	cfg := &apiConfig{}
	payload := map[string]string{"body": "hello world"}
	b, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/validate_chirp", bytes.NewReader(b))
	rr := httptest.NewRecorder()

	cfg.validateChirpHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}
	var resp map[string]string
	json.Unmarshal(rr.Body.Bytes(), &resp)
	expected := "hello world"
	if resp["cleaned_body"] != expected {
		t.Errorf("expected cleaned_body %q, got %q", expected, resp["cleaned_body"])
	}
}

func TestValidateChirpHandler_EmptyBody(t *testing.T) {
	cfg := &apiConfig{}
	payload := map[string]string{"body": ""}
	b, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/validate_chirp", bytes.NewReader(b))
	rr := httptest.NewRecorder()

	cfg.validateChirpHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}
	body, _ := io.ReadAll(rr.Body)
	var resp map[string]string
	json.Unmarshal(body, &resp)
	if resp["cleaned_body"] != "" {
		t.Errorf("expected empty cleaned_body, got %q", resp["cleaned_body"])
	}
}
