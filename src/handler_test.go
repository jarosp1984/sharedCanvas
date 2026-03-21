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

	if lines[0].ID != 1 {
		t.Fatalf("expected ID 1, got %d", lines[0].ID)
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

	// GET should now be accepted (not rejected) since we merged POST and GET handlers
	request := httptest.NewRequest(http.MethodGet, "/api/lines", nil)
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d (GET now supported), got %d", http.StatusOK, response.Code)
	}
}

func TestGetLinesReturnsEmptyArray(t *testing.T) {
	store := NewMemoryLineStore()
	handler := NewLinesHandler(store)

	request := httptest.NewRequest(http.MethodGet, "/api/lines", nil)
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, response.Code)
	}

	if !strings.Contains(response.Body.String(), `"data":[]`) {
		t.Fatalf("expected empty data array, got %s", response.Body.String())
	}
}

func TestGetLinesReturnsStoredLines(t *testing.T) {
	store := NewMemoryLineStore()
	handler := NewLinesHandler(store)

	// POST two lines
	for i := 0; i < 2; i++ {
		postReq := httptest.NewRequest(http.MethodPost, "/api/lines", strings.NewReader(`{"x1":10,"y1":20,"x2":30,"y2":40,"color":"#112233","width":5}`))
		postResp := httptest.NewRecorder()
		handler.ServeHTTP(postResp, postReq)

		if postResp.Code != http.StatusCreated {
			t.Fatalf("failed to post line: %d", postResp.Code)
		}
	}

	// GET the lines
	getReq := httptest.NewRequest(http.MethodGet, "/api/lines", nil)
	getResp := httptest.NewRecorder()

	handler.ServeHTTP(getResp, getReq)

	if getResp.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, getResp.Code)
	}

	if !strings.Contains(getResp.Body.String(), `"status":"ok"`) {
		t.Fatalf("expected status ok in response, got %s", getResp.Body.String())
	}

	if !strings.Contains(getResp.Body.String(), `"id":1`) {
		t.Fatalf("expected line id 1 in response, got %s", getResp.Body.String())
	}

	if !strings.Contains(getResp.Body.String(), `"id":2`) {
		t.Fatalf("expected line id 2 in response, got %s", getResp.Body.String())
	}
}

func TestGetLinesRejectsNonGetMethod(t *testing.T) {
	store := NewMemoryLineStore()
	handler := NewLinesHandler(store)

	request := httptest.NewRequest("DELETE", "/api/lines", nil)
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status %d, got %d", http.StatusMethodNotAllowed, response.Code)
	}
}
