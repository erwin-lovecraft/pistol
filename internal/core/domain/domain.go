package domain

import (
	"encoding/json"
	"net/http"
)

type Room struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

type Event struct {
	Method      string              `json:"method"`
	Header      http.Header         `json:"header"`
	QueryParams map[string][]string `json:"query_params"`
	Body        json.RawMessage     `json:"body,omitempty"`
}
