package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type APIFunc func(http.ResponseWriter, *http.Request) error

type APIServer struct {
	listenAddr string
	store      Storage
}

type APIError struct {
	Error string
}

func NewApiServer(listenAddr string) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
	}
}

func (s *APIServer) Run() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /account", makeHTTPHandleFunc(s.handleGetAccount))
	mux.HandleFunc("GET /account/{id}", makeHTTPHandleFunc(s.handleGetAccountById))
	mux.HandleFunc("POST /account", makeHTTPHandleFunc(s.handleCreateAccount))
	mux.HandleFunc("DELETE /account", makeHTTPHandleFunc(s.handleDeleteAccount))

	log.Println("API server running on port: ", s.listenAddr)

	http.ListenAndServe(s.listenAddr, mux)
}

func (s *APIServer) handleGetAccount(w http.ResponseWriter, req *http.Request) error {
	account := NewAccount("Cathal", "OC")

	return WriteJSON(w, http.StatusOK, account)
}

func (s *APIServer) handleGetAccountById(w http.ResponseWriter, req *http.Request) error {
	id := req.PathValue("id")

	account := fmt.Sprintf("Got account for id: %s", id)

	return WriteJSON(w, http.StatusOK, account)
}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, req *http.Request) error {
	return nil
}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, req *http.Request) error {
	return nil
}

func (s *APIServer) handleTransfer(w http.ResponseWriter, req *http.Request) error {
	return nil
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func makeHTTPHandleFunc(f APIFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if err := f(w, req); err != nil {
			//handle error here
			WriteJSON(w, http.StatusBadRequest, APIError{Error: err.Error()})
		}
	}
}
