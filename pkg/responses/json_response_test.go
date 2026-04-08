package responses

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestJsonResponseStatusCode(t *testing.T) {
	resp := JsonResponse(http.StatusOK, map[string]string{"key": "value"})
	if resp.StatusCode != http.StatusOK {
		t.Errorf("StatusCode = %d; want %d", resp.StatusCode, http.StatusOK)
	}
}

func TestJsonResponseContentTypeHeader(t *testing.T) {
	resp := JsonResponse(http.StatusCreated, nil)
	if resp.Headers["Content-Type"] != "application/json" {
		t.Errorf("Content-Type = %q; want application/json", resp.Headers["Content-Type"])
	}
}

func TestJsonResponseBody(t *testing.T) {
	body := map[string]string{"hello": "world"}
	resp := JsonResponse(http.StatusOK, body)
	var got map[string]string
	if err := json.Unmarshal([]byte(resp.Body), &got); err != nil {
		t.Fatalf("failed to unmarshal body: %v", err)
	}
	if got["hello"] != "world" {
		t.Errorf("body hello = %q; want world", got["hello"])
	}
}

func TestJsonResponse201(t *testing.T) {
	resp := JsonResponse(http.StatusCreated, map[string]string{"id": "123"})
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("StatusCode = %d; want %d", resp.StatusCode, http.StatusCreated)
	}
}

func TestJsonResponse404(t *testing.T) {
	resp := JsonResponse(http.StatusNotFound, map[string]string{"error": "not found"})
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("StatusCode = %d; want %d", resp.StatusCode, http.StatusNotFound)
	}
	var got map[string]string
	if err := json.Unmarshal([]byte(resp.Body), &got); err != nil {
		t.Fatalf("failed to unmarshal body: %v", err)
	}
	if got["error"] != "not found" {
		t.Errorf("error = %q; want not found", got["error"])
	}
}

func TestJsonResponseNilBody(t *testing.T) {
	resp := JsonResponse(http.StatusNoContent, nil)
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("StatusCode = %d; want %d", resp.StatusCode, http.StatusNoContent)
	}
	if resp.Body != "null" {
		t.Errorf("Body = %q; want null for nil body", resp.Body)
	}
}
