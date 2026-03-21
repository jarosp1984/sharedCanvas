package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateLineHandlerStoresValidLine(t *testing.T) {
	store := NewMemoryLineStore()
	handler := NewCreateLineHandler(store)

	request := httptest.NewRequest(http.MethodPost, "/api/lines", strings.NewReader(`{"x1":10,"y1":20,"x2":30,"y2":40,"color":"#112233","width":5}`))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, response.Code)
	}

	lines, err := store.ListLines(context.Background())
	if err != nil {
		t.Fatalf("list lines returned error: %v", err)
	}

	if len(lines) != 1 {
		t.Fatalf("expected 1 line, got %d", len(lines))
	}

	if lines[0].Color != "#112233" {
		t.Fatalf("expected stored color #112233, got %s", lines[0].Color)
	}
}

func TestCreateLineHandlerRejectsInvalidLine(t *testing.T) {
	store := NewMemoryLineStore()
	handler := NewCreateLineHandler(store)

	request := httptest.NewRequest(http.MethodPost, "/api/lines", strings.NewReader(`{"x1":-1,"y1":20,"x2":30,"y2":40,"color":"#112233","width":5}`))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, response.Code)
	}
}

func TestCreateLineHandlerRejectsWrongMethod(t *testing.T) {
	store := NewMemoryLineStore()
	handler := NewCreateLineHandler(store)

	request := httptest.NewRequest(http.MethodGet, "/api/lines", nil)
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status %d, got %d", http.StatusMethodNotAllowed, response.Code)
	}
}
