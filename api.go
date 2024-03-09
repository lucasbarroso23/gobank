package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// The json.NewEnconder needs the io.Writer, and the http.ResponseWriter implement it
// luckily this way we can encode whatever ´v´ is, where ´v´is equivalent to an empty interface
func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}

// basically this type is used to be able to conform our handle methods to http.HandlerFunc interface
// the http.HandlerFunc by itself does not return an error and since our method does, we add this layer
// to make this work with our mux router.HandleFunc
type apiFunc func(http.ResponseWriter, *http.Request) error

type apiError struct {
	Error string
}

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, apiError{Error: err.Error()})
		}
	}
}

type APIServer struct {
	listenAddr string
}

func NewApiServer(listenAddr string) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
	}
}

// We use mux because even tho the http package works pretty well to handle the requests
// there is no helping methods to handle the path/params/type
// TODO: checkout the updates to http/net packages, i beliebe that some improvements were
// made related to handling the path and params
func (s *APIServer) run() {
	router := mux.NewRouter()

	router.HandleFunc("/account", makeHTTPHandleFunc(s.handleAccount))
	router.HandleFunc("/account/{id}", makeHTTPHandleFunc(s.handleGetAccount))

	log.Println("JSON API server running on port: ", s.listenAddr)

	http.ListenAndServe(s.listenAddr, router)
}

// since we cant specified the request type using mux, this method is used to check the
// request type and send to the correct handler
func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	// TODO: try using enum and a switch case to handle the requests methods and routing
	if r.Method == "GET" {
		return s.handleGetAccount(w, r)
	}
	if r.Method == "POST" {
		return s.handleCreateAccount(w, r)
	}
	if r.Method == "DELETE" {
		return s.handleDeleteAccount(w, r)
	}

	return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	id := mux.Vars(r)["id"]

	fmt.Println(id)

	return WriteJSON(w, http.StatusOK, &Account{})
}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	return nil
}
