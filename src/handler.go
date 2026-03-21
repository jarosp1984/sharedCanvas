package main

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
)

func NewCreateLineHandler(store LineStore) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodPost {
			writeJSONError(writer, http.StatusMethodNotAllowed, "method must be POST")
			return
		}

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

		if err := store.SaveLine(request.Context(), line); err != nil {
			log.Printf("save line failed: %v", err)
			writeJSONError(writer, http.StatusInternalServerError, "failed to save line")
			return
		}

		writeJSON(writer, http.StatusCreated, map[string]any{
			"status": "ok",
			"line":   line,
		})
	}
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
