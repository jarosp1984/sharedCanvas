package main

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"
)

func NewLinesHandler(store LineStore) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		sessionID, err := sessionIDFromAPIPath(request.URL.Path, linesPathSuffix)
		if err != nil {
			handleSessionRouteError(writer, err)
			return
		}

		switch request.Method {
		case http.MethodPost:
			handlePostLine(writer, request, store, sessionID)
		case http.MethodGet:
			handleGetLines(writer, request, store, sessionID)
		default:
			writeJSONError(writer, http.StatusMethodNotAllowed, "method must be GET or POST")
		}
	}
}

func NewClearLinesHandler(store LineStore) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		sessionID, err := sessionIDFromAPIPath(request.URL.Path, clearLinesPathSuffix)
		if err != nil {
			handleSessionRouteError(writer, err)
			return
		}

		if request.Method != http.MethodPost {
			writeJSONError(writer, http.StatusMethodNotAllowed, "method must be POST")
			return
		}

		handleClearLines(writer, request, store, sessionID)
	}
}

func NewSessionAPIHandler(store LineStore) http.HandlerFunc {
	linesHandler := NewLinesHandler(store)
	clearHandler := NewClearLinesHandler(store)

	return func(writer http.ResponseWriter, request *http.Request) {
		request.URL.Path = normalizeSessionPath(request.URL.Path)

		switch {
		case hasPathSuffix(request.URL.Path, clearLinesPathSuffix):
			clearHandler.ServeHTTP(writer, request)
		case hasPathSuffix(request.URL.Path, linesPathSuffix):
			linesHandler.ServeHTTP(writer, request)
		default:
			http.NotFound(writer, request)
		}
	}
}

func NewRootHandler() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/" {
			http.NotFound(writer, request)
			return
		}

		sessionID, err := newSessionID()
		if err != nil {
			log.Printf("create session failed: %v", err)
			writeJSONError(writer, http.StatusInternalServerError, "failed to create session")
			return
		}

		http.Redirect(writer, request, sessionPagePathPrefix+sessionID, http.StatusSeeOther)
	}
}

func NewSessionPageHandler() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		_, err := sessionIDFromPagePath(request.URL.Path)
		if err != nil {
			if errors.Is(err, errInvalidSessionID) {
				http.NotFound(writer, request)
				return
			}

			http.NotFound(writer, request)
			return
		}

		http.ServeFile(writer, request, staticIndexPath)
	}
}

func handlePostLine(writer http.ResponseWriter, request *http.Request, store LineStore, sessionID string) {
	request.Body = http.MaxBytesReader(writer, request.Body, 1<<20)

	var line Line
	decoder := json.NewDecoder(request.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&line); err != nil {
		writeJSONError(writer, http.StatusBadRequest, "invalid JSON payload")
		return
	}

	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		writeJSONError(writer, http.StatusBadRequest, "request body must contain a single JSON object")
		return
	}

	if err := line.Validate(); err != nil {
		writeJSONError(writer, http.StatusBadRequest, err.Error())
		return
	}

	if err := store.SaveLine(request.Context(), sessionID, &line); err != nil {
		log.Printf("save line failed: %v", err)
		writeJSONError(writer, http.StatusInternalServerError, "failed to save line")
		return
	}

	writeJSON(writer, http.StatusCreated, map[string]any{
		"status": "ok",
		"line":   line,
	})
}

func handleGetLines(writer http.ResponseWriter, request *http.Request, store LineStore, sessionID string) {
	sinceValue := request.URL.Query().Get("since")

	var (
		lines []Line
		err   error
	)

	if sinceValue == "" {
		lines, err = store.ListLines(request.Context(), sessionID)
	} else {
		sinceID, parseErr := strconv.Atoi(sinceValue)
		if parseErr != nil || sinceID < 0 {
			writeJSONError(writer, http.StatusBadRequest, "since must be a non-negative integer")
			return
		}

		lines, err = store.ListLinesSince(request.Context(), sessionID, sinceID)
	}

	if err != nil {
		log.Printf("list lines failed: %v", err)
		writeJSONError(writer, http.StatusInternalServerError, "failed to list lines")
		return
	}

	writeJSON(writer, http.StatusOK, map[string]any{
		"status": "ok",
		"data":   lines,
	})
}

func handleClearLines(writer http.ResponseWriter, request *http.Request, store LineStore, sessionID string) {
	if err := store.ClearLines(request.Context(), sessionID); err != nil {
		log.Printf("clear lines failed: %v", err)
		writeJSONError(writer, http.StatusInternalServerError, "failed to clear lines")
		return
	}

	writer.WriteHeader(http.StatusNoContent)
}

func handleSessionRouteError(writer http.ResponseWriter, err error) {
	if errors.Is(err, errInvalidSessionID) {
		writeJSONError(writer, http.StatusBadRequest, err.Error())
		return
	}

	writeJSONError(writer, http.StatusNotFound, "not found")
}

func hasPathSuffix(path string, suffix string) bool {
	return len(path) >= len(suffix) && path[len(path)-len(suffix):] == suffix
}

// Deprecated: use NewSessionAPIHandler instead
func NewCreateLineHandler(store LineStore) http.HandlerFunc {
	return NewSessionAPIHandler(store)
}

func writeJSON(writer http.ResponseWriter, status int, payload any) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)

	if err := json.NewEncoder(writer).Encode(payload); err != nil {
		log.Printf("write JSON failed: %v", err)
	}
}

func writeJSONError(writer http.ResponseWriter, status int, message string) {
	writeJSON(writer, status, map[string]string{"error": message})
}
