package service

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRootHandler(t *testing.T) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(RootHandler)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %+v want %+v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := "{\"errors\":[{\"error\":\"cannot handle request\"},{\"error\":\"not found\"}]}"
	if rr.Body.String() != expected {
		fmt.Printf("handler returned unexpected body: got %v want %v",
			rr.Body.Bytes(), []byte(expected))
		t.Errorf("handler returned unexpected body: got %s want %s",
			rr.Body.Bytes(), []byte(expected))
	}
}
