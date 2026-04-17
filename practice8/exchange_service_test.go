package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetRate_SuccessfulScenario(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/convert", r.URL.Path)
		assert.Equal(t, "USD", r.URL.Query().Get("from"))
		assert.Equal(t, "EUR", r.URL.Query().Get("to"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"base": "USD", "target": "EUR", "rate": 0.85}`))
	}))
	defer server.Close()

	service := NewExchangeService(server.URL)
	rate, err := service.GetRate("USD", "EUR")

	assert.NoError(t, err)
	assert.Equal(t, 0.85, rate)
}

func TestGetRate_APIBusinessError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "invalid currency pair"}`))
	}))
	defer server.Close()

	service := NewExchangeService(server.URL)
	rate, err := service.GetRate("USD", "UNKNOWN")

	assert.Error(t, err)
	assert.Equal(t, 0.0, rate)
	assert.Contains(t, err.Error(), "api error: invalid currency pair")
}

func TestGetRate_MalformedJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{invalid json`))
	}))
	defer server.Close()

	service := NewExchangeService(server.URL)
	rate, err := service.GetRate("USD", "EUR")

	assert.Error(t, err)
	assert.Equal(t, 0.0, rate)
	assert.Contains(t, err.Error(), "decode error")
}

func TestGetRate_SlowResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(6 * time.Second) // Timeout is 5 seconds
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"base": "USD", "target": "EUR", "rate": 0.85}`))
	}))
	defer server.Close()

	service := NewExchangeService(server.URL)
	
	rate, err := service.GetRate("USD", "EUR")

	assert.Error(t, err)
	assert.Equal(t, 0.0, rate)
	assert.Contains(t, err.Error(), "network error")
	assert.True(t, strings.Contains(err.Error(), "context deadline exceeded") || strings.Contains(err.Error(), "Timeout"), "error should indicate timeout")
}

func TestGetRate_ServerPanic(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "internal server error"}`))
	}))
	defer server.Close()

	service := NewExchangeService(server.URL)
	rate, err := service.GetRate("USD", "EUR")

	assert.Error(t, err)
	assert.Equal(t, 0.0, rate)
	assert.Contains(t, err.Error(), "internal server error")
}

func TestGetRate_ServerPanic_NoJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`Internal Server Error`))
	}))
	defer server.Close()

	service := NewExchangeService(server.URL)
	rate, err := service.GetRate("USD", "EUR")

	assert.Error(t, err)
	assert.Equal(t, 0.0, rate)
	assert.Contains(t, err.Error(), "decode error")
}

func TestGetRate_EmptyBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	service := NewExchangeService(server.URL)
	rate, err := service.GetRate("USD", "EUR")

	assert.Error(t, err)
	assert.Equal(t, 0.0, rate)
	assert.Contains(t, err.Error(), "decode error")
}
