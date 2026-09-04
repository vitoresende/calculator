package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBuildServer_DefaultAndCustomConfig(t *testing.T) {
	t.Run("default configuration", func(t *testing.T) {
		server := buildServer("", "")
		if server == nil {
			t.Fatal("expected server not to be nil")
		}
		if server.Addr != ":8080" {
			t.Errorf("expected default port :8080, got %s", server.Addr)
		}
		if server.Handler == nil {
			t.Fatal("expected server handler to be configured")
		}

		// Verify health endpoint on configured handler
		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		rec := httptest.NewRecorder()
		server.Handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected status 200 for health route on server, got %d", rec.Code)
		}
	})

	t.Run("custom port and origin", func(t *testing.T) {
		server := buildServer("9090", "https://custom.domain.com")
		if server == nil {
			t.Fatal("expected server not to be nil")
		}
		if server.Addr != ":9090" {
			t.Errorf("expected port :9090, got %s", server.Addr)
		}
	})
}
