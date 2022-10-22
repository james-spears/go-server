package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

type message struct {
	Message string `json:"message"`
}

// DefaultHandler is the generic 404 response.
func DefaultHandler(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	_, err := io.WriteString(w, "Not Found.")
	if err != nil {
		log.Print(err)
	}
}

type HelloWorldRequest struct {
	Name string `json:"name"`
}

type HelloWorldResponse struct {
	Message string `json:"name"`
}

// HelloWorldHandler wraps the geocoding functionality provided
// by the go client library.
func HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	var request HelloWorldRequest
	var response HelloWorldResponse
	// SET THE GENERIC CONTENT-TYPE HEADER TO JSON
	encoder := json.NewEncoder(w)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	switch method := r.Method; method {
	case http.MethodGet:
		response = HelloWorldResponse{Message: "Hello, World!"}
		w.WriteHeader(http.StatusOK)
		err := encoder.Encode(response)
		if err != nil {
			log.Fatal(err)
		}
	case http.MethodPost:
		b, err := io.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			err = r.Body.Close()
			if err != nil {
				log.Fatal(err)
			}
		}()

		err = json.Unmarshal(b, &request)
		if err != nil {
			// If the request body cannot be unmarshalled
			// into the QueryAutocompleteRequest struct
			// then it is very likely due to an improperly
			// formatted request body.
			w.WriteHeader(http.StatusBadRequest)
			err = encoder.Encode(message{Message: "bad request"})
			if err != nil {
				log.Fatal(err)
			}
		}
		response = HelloWorldResponse{Message: fmt.Sprintf("Hello, %s!", request.Name)}
		w.WriteHeader(http.StatusOK)
		err = encoder.Encode(response)
		if err != nil {
			log.Fatal(err)
		}
	case http.MethodOptions:
		// We must allow options method for preflight requests.
		w.WriteHeader(http.StatusOK)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		err := encoder.Encode(message{Message: "method not allowed"})
		if err != nil {
			log.Fatal(err)
		}
	}
}

// V1Handler is a simple switch routing the request to
// based on the version container in the url.
func V1Handler(w http.ResponseWriter, r *http.Request) {
	url := strings.Split(r.URL.Path, "/")
	switch service := url[3]; service {
	// Add service handlers here ...
	case "hello-world":
		HelloWorldHandler(w, r)
		return
	default:
		DefaultHandler(w)
		return
	}
}

// ApiHandler is a simple switch routing the request to
// based on the version (after base) url.
func ApiHandler(w http.ResponseWriter, r *http.Request) {
	url := strings.Split(r.URL.Path, "/")
	switch version := url[2]; version {
	case "v1":
		V1Handler(w, r)
	default:
		DefaultHandler(w)
		return
	}
}

// RootHandler is a simple switch routing the request to
// based on the first URL component.
func RootHandler(w http.ResponseWriter, r *http.Request) {
	// Do not require trailing slash. Add one if not present.
	if string(r.URL.Path[len(r.URL.Path)-1]) != "/" {
		r.URL.Path += "/"
	}
	url := strings.Split(r.URL.Path, "/")
	switch base := url[1]; base {
	case "api":
		w.Header().Set("Content-Type", "application/json")
		ApiHandler(w, r)
	default:
		w.Header().Set("Content-Type", "application/plain")
		DefaultHandler(w)
		return
	}
}
