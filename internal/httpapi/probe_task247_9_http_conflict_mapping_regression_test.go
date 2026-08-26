package httpapi

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"task247-reactcons/internal/service"
	"task247-reactcons/internal/store"
)

func TestTask247Bug09HTTPConflictMapping(t *testing.T) {
	repository, err := store.Open(t.TempDir() + "/reactcons.db")
	if err != nil {
		t.Fatal(err)
	}
	defer repository.Close()
	handler := New(service.New(repository)).Handler()
	payload := []byte(`{"ext_key":"same-network","name":"first"}`)
	first := httptest.NewRecorder()
	handler.ServeHTTP(first, httptest.NewRequest(http.MethodPost, "/api/networks", bytes.NewReader(payload)))
	if first.Code != http.StatusCreated {
		t.Fatalf("first create failed: %d %s", first.Code, first.Body.String())
	}
	second := httptest.NewRecorder()
	handler.ServeHTTP(second, httptest.NewRequest(http.MethodPost, "/api/networks", bytes.NewReader(payload)))
	if second.Code != http.StatusConflict || !bytes.Contains(second.Body.Bytes(), []byte("conflict")) {
		t.Fatalf("duplicate create was not a conflict response: %d %s", second.Code, second.Body.String())
	}
}
