package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
)

const (
	testSessionID      = "11111111-1111-4111-8111-111111111111"
	otherTestSessionID = "22222222-2222-4222-8222-222222222222"
)

func newTestServer() (*MemoryLineStore, http.Handler) {
	store := NewMemoryLineStore()

	return store, NewServerMux(store)
}

func performRequest(handler http.Handler, method string, path string, body string) *httptest.ResponseRecorder {
	request := httptest.NewRequest(method, path, strings.NewReader(body))
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)

	return response
}

func linesPath(sessionID string) string {
	return "/api/sessions/" + sessionID + "/lines"
}

func clearLinesPath(sessionID string) string {
	return "/api/sessions/" + sessionID + "/lines/clear"
}

func validLineJSON() string {
	return `{"x1":10,"y1":20,"x2":30,"y2":40,"color":"#112233","width":5}`
}

func TestSessionLinesHandlerStoresValidLine(t *testing.T) {
	store, handler := newTestServer()

	response := performRequest(handler, http.MethodPost, linesPath(testSessionID), validLineJSON())
	if response.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, response.Code)
	}

	lines, err := store.ListLines(context.Background(), testSessionID)
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

func TestSessionLinesHandlerRejectsInvalidLine(t *testing.T) {
	_, handler := newTestServer()

	response := performRequest(handler, http.MethodPost, linesPath(testSessionID), `{"x1":-1,"y1":20,"x2":30,"y2":40,"color":"#112233","width":5}`)
	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, response.Code)
	}
}

func TestSessionLinesHandlerRejectsInvalidSessionID(t *testing.T) {
	_, handler := newTestServer()

	response := performRequest(handler, http.MethodGet, linesPath("not-a-uuid"), "")
	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, response.Code)
	}

	if !strings.Contains(response.Body.String(), `"error":"session_id must be a valid UUID"`) {
		t.Fatalf("expected invalid session id error, got %s", response.Body.String())
	}
}

func TestSessionLinesHandlerRejectsWrongMethod(t *testing.T) {
	_, handler := newTestServer()

	response := performRequest(handler, http.MethodDelete, linesPath(testSessionID), "")
	if response.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status %d, got %d", http.StatusMethodNotAllowed, response.Code)
	}
}

func TestGetLinesReturnsEmptyArrayPerSession(t *testing.T) {
	_, handler := newTestServer()

	response := performRequest(handler, http.MethodGet, linesPath(testSessionID), "")
	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, response.Code)
	}

	if !strings.Contains(response.Body.String(), `"data":[]`) {
		t.Fatalf("expected empty data array, got %s", response.Body.String())
	}
}

func TestSessionLinesAreIsolatedBetweenSessions(t *testing.T) {
	_, handler := newTestServer()

	postResponse := performRequest(handler, http.MethodPost, linesPath(testSessionID), validLineJSON())
	if postResponse.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, postResponse.Code)
	}

	getOtherResponse := performRequest(handler, http.MethodGet, linesPath(otherTestSessionID), "")
	if getOtherResponse.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, getOtherResponse.Code)
	}

	if !strings.Contains(getOtherResponse.Body.String(), `"data":[]`) {
		t.Fatalf("expected other session to be empty, got %s", getOtherResponse.Body.String())
	}
}

func TestGetLinesReturnsStoredLinesForSession(t *testing.T) {
	_, handler := newTestServer()

	for i := 0; i < 2; i++ {
		postResponse := performRequest(handler, http.MethodPost, linesPath(testSessionID), validLineJSON())
		if postResponse.Code != http.StatusCreated {
			t.Fatalf("failed to post line: %d", postResponse.Code)
		}
	}

	getResponse := performRequest(handler, http.MethodGet, linesPath(testSessionID), "")
	if getResponse.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, getResponse.Code)
	}

	if !strings.Contains(getResponse.Body.String(), `"id":1`) {
		t.Fatalf("expected line id 1 in response, got %s", getResponse.Body.String())
	}

	if !strings.Contains(getResponse.Body.String(), `"id":2`) {
		t.Fatalf("expected line id 2 in response, got %s", getResponse.Body.String())
	}
}

func TestGetLinesSinceReturnsOnlyNewerLinesWithinSession(t *testing.T) {
	_, handler := newTestServer()

	for i := 0; i < 3; i++ {
		postResponse := performRequest(handler, http.MethodPost, linesPath(testSessionID), validLineJSON())
		if postResponse.Code != http.StatusCreated {
			t.Fatalf("failed to post line: %d", postResponse.Code)
		}
	}

	otherPostResponse := performRequest(handler, http.MethodPost, linesPath(otherTestSessionID), validLineJSON())
	if otherPostResponse.Code != http.StatusCreated {
		t.Fatalf("failed to post line to other session: %d", otherPostResponse.Code)
	}

	getResponse := performRequest(handler, http.MethodGet, linesPath(testSessionID)+"?since=1", "")
	if getResponse.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, getResponse.Code)
	}

	body := getResponse.Body.String()
	if strings.Contains(body, `"id":1`) {
		t.Fatalf("expected id 1 to be filtered out, got %s", body)
	}

	if !strings.Contains(body, `"id":2`) || !strings.Contains(body, `"id":3`) {
		t.Fatalf("expected ids 2 and 3 in response, got %s", body)
	}
}

func TestGetLinesSinceRejectsInvalidValue(t *testing.T) {
	_, handler := newTestServer()

	response := performRequest(handler, http.MethodGet, linesPath(testSessionID)+"?since=-1", "")
	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, response.Code)
	}

	if !strings.Contains(response.Body.String(), `"error":"since must be a non-negative integer"`) {
		t.Fatalf("expected validation error in response, got %s", response.Body.String())
	}
}

func TestClearLinesRemovesOnlyTargetSessionLines(t *testing.T) {
	_, handler := newTestServer()

	for _, sessionID := range []string{testSessionID, otherTestSessionID} {
		for i := 0; i < 2; i++ {
			postResponse := performRequest(handler, http.MethodPost, linesPath(sessionID), validLineJSON())
			if postResponse.Code != http.StatusCreated {
				t.Fatalf("failed to post line: %d", postResponse.Code)
			}
		}
	}

	clearResponse := performRequest(handler, http.MethodPost, clearLinesPath(testSessionID), "")
	if clearResponse.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, clearResponse.Code)
	}

	getClearedResponse := performRequest(handler, http.MethodGet, linesPath(testSessionID), "")
	if !strings.Contains(getClearedResponse.Body.String(), `"data":[]`) {
		t.Fatalf("expected cleared session to be empty, got %s", getClearedResponse.Body.String())
	}

	getOtherResponse := performRequest(handler, http.MethodGet, linesPath(otherTestSessionID), "")
	if !strings.Contains(getOtherResponse.Body.String(), `"id":1`) || !strings.Contains(getOtherResponse.Body.String(), `"id":2`) {
		t.Fatalf("expected other session lines to remain, got %s", getOtherResponse.Body.String())
	}
}

func TestClearLinesRejectsNonPostMethod(t *testing.T) {
	_, handler := newTestServer()

	response := performRequest(handler, http.MethodGet, clearLinesPath(testSessionID), "")
	if response.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status %d, got %d", http.StatusMethodNotAllowed, response.Code)
	}
}

func TestClearLinesPreservesMonotonicIDsPerSession(t *testing.T) {
	_, handler := newTestServer()

	firstPostResponse := performRequest(handler, http.MethodPost, linesPath(testSessionID), validLineJSON())
	if firstPostResponse.Code != http.StatusCreated {
		t.Fatalf("failed to post first line: %d", firstPostResponse.Code)
	}

	clearResponse := performRequest(handler, http.MethodPost, clearLinesPath(testSessionID), "")
	if clearResponse.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, clearResponse.Code)
	}

	secondPostResponse := performRequest(handler, http.MethodPost, linesPath(testSessionID), `{"x1":15,"y1":25,"x2":35,"y2":45,"color":"#112233","width":5}`)
	if secondPostResponse.Code != http.StatusCreated {
		t.Fatalf("failed to post second line: %d", secondPostResponse.Code)
	}

	if !strings.Contains(secondPostResponse.Body.String(), `"id":2`) {
		t.Fatalf("expected second line id to be 2 after clear, got %s", secondPostResponse.Body.String())
	}
}

func TestLineIDsStartAtOnePerSession(t *testing.T) {
	_, handler := newTestServer()

	firstSessionResponse := performRequest(handler, http.MethodPost, linesPath(testSessionID), validLineJSON())
	if !strings.Contains(firstSessionResponse.Body.String(), `"id":1`) {
		t.Fatalf("expected first session to start at id 1, got %s", firstSessionResponse.Body.String())
	}

	secondSessionResponse := performRequest(handler, http.MethodPost, linesPath(otherTestSessionID), validLineJSON())
	if !strings.Contains(secondSessionResponse.Body.String(), `"id":1`) {
		t.Fatalf("expected second session to start at id 1, got %s", secondSessionResponse.Body.String())
	}
}

func TestRootRedirectsToNewSession(t *testing.T) {
	_, handler := newTestServer()

	response := performRequest(handler, http.MethodGet, "/", "")
	if response.Code != http.StatusSeeOther {
		t.Fatalf("expected status %d, got %d", http.StatusSeeOther, response.Code)
	}

	location := response.Header().Get("Location")
	pattern := regexp.MustCompile(`^/sessions/[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)
	if !pattern.MatchString(location) {
		t.Fatalf("expected redirect to session path, got %q", location)
	}
}

func TestSessionPageServesCanvasForValidSession(t *testing.T) {
	_, handler := newTestServer()

	response := performRequest(handler, http.MethodGet, "/sessions/"+testSessionID, "")
	if response.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, response.Code)
	}

	if !strings.Contains(response.Body.String(), "Shared Canvas - Paint Tool") {
		t.Fatalf("expected HTML canvas page, got %s", response.Body.String())
	}
}

func TestSessionPageRejectsInvalidSessionID(t *testing.T) {
	_, handler := newTestServer()

	response := performRequest(handler, http.MethodGet, "/sessions/not-a-uuid", "")
	if response.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, response.Code)
	}
}
