package main

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"regexp"
	"strings"
)

const (
	sessionPagePathPrefix = "/sessions/"
	apiSessionsPathPrefix = "/api/sessions/"
	linesPathSuffix       = "/lines"
	clearLinesPathSuffix  = "/lines/clear"
	staticIndexPath       = "../static/index.html"
)

var (
	errSessionRouteNotFound = errors.New("session route not found")
	errInvalidSessionID     = errors.New("session_id must be a valid UUID")
	sessionIDPattern        = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-4[0-9a-fA-F]{3}-[89abAB][0-9a-fA-F]{3}-[0-9a-fA-F]{12}$`)
)

func newSessionID() (string, error) {
	var raw [16]byte
	if _, err := rand.Read(raw[:]); err != nil {
		return "", err
	}

	raw[6] = (raw[6] & 0x0f) | 0x40
	raw[8] = (raw[8] & 0x3f) | 0x80

	var encoded [36]byte
	hex.Encode(encoded[0:8], raw[0:4])
	encoded[8] = '-'
	hex.Encode(encoded[9:13], raw[4:6])
	encoded[13] = '-'
	hex.Encode(encoded[14:18], raw[6:8])
	encoded[18] = '-'
	hex.Encode(encoded[19:23], raw[8:10])
	encoded[23] = '-'
	hex.Encode(encoded[24:36], raw[10:16])

	return string(encoded[:]), nil
}

func validateSessionID(sessionID string) error {
	if !sessionIDPattern.MatchString(sessionID) {
		return errInvalidSessionID
	}

	return nil
}

func normalizeSessionPath(path string) string {
	if len(path) > 1 {
		return strings.TrimSuffix(path, "/")
	}

	return path
}

func sessionIDFromPagePath(path string) (string, error) {
	normalizedPath := normalizeSessionPath(path)
	if !strings.HasPrefix(normalizedPath, sessionPagePathPrefix) {
		return "", errSessionRouteNotFound
	}

	sessionID := strings.TrimPrefix(normalizedPath, sessionPagePathPrefix)
	if sessionID == "" || strings.Contains(sessionID, "/") {
		return "", errSessionRouteNotFound
	}

	if err := validateSessionID(sessionID); err != nil {
		return "", err
	}

	return sessionID, nil
}

func sessionIDFromAPIPath(path string, suffix string) (string, error) {
	normalizedPath := normalizeSessionPath(path)
	if !strings.HasPrefix(normalizedPath, apiSessionsPathPrefix) || !strings.HasSuffix(normalizedPath, suffix) {
		return "", errSessionRouteNotFound
	}

	sessionID := strings.TrimSuffix(strings.TrimPrefix(normalizedPath, apiSessionsPathPrefix), suffix)
	if sessionID == "" || strings.Contains(sessionID, "/") {
		return "", errSessionRouteNotFound
	}

	if err := validateSessionID(sessionID); err != nil {
		return "", err
	}

	return sessionID, nil
}
