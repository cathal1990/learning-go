package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type APIFunc func(http.ResponseWriter, *http.Request) error

type APIServer struct {
	listenAddr string
	store      Storage
}

type APIError struct {
	Error string `json:"error"`
}

func NewApiServer(listenAddr string, store Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

func (s *APIServer) Run() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /accounts", makeHTTPHandleFunc(s.handleGetAccounts))
	mux.HandleFunc("GET /account", makeHTTPHandleFunc(s.handleGetAccount))
	mux.HandleFunc("GET /account/{id}", makeHTTPHandleFunc(s.handleGetAccountById))
	mux.HandleFunc("POST /account", makeHTTPHandleFunc(s.handleCreateAccount))
	mux.HandleFunc("DELETE /account/{id}", makeHTTPHandleFunc(s.handleDeleteAccount))

	log.Println("API server running on port: ", s.listenAddr)

	http.ListenAndServe(s.listenAddr, mux)
}

func (s *APIServer) handleGetAccount(w http.ResponseWriter, req *http.Request) error {
	account := NewAccount("Cathal", "OC")

	return WriteJSON(w, http.StatusOK, account)
}

func (s *APIServer) handleGetAccountById(w http.ResponseWriter, req *http.Request) error {
	id, err := getIdFromParams(req)

	if err != nil {
		return WriteJSON(w, http.StatusBadRequest, err)
	}

	account, err := s.store.GetAccountById(id)

	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, account)
}

func (s *APIServer) handleGetAccounts(w http.ResponseWriter, req *http.Request) error {
	accounts, err := s.store.GetAccounts()

	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, accounts)
}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, req *http.Request) error {
	createAccountReq := new(CreateAccountRequest)

	if err := json.NewDecoder(req.Body).Decode(createAccountReq); err != nil {
		return err
	}

	account := NewAccount(createAccountReq.FirstName, createAccountReq.LastName)

	if err := s.store.CreateAccount(account); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, account)
}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, req *http.Request) error {
	id, err := getIdFromParams(req)

	if err != nil {
		return WriteJSON(w, http.StatusBadRequest, err)
	}

	deleteErr := s.store.DeleteAccount(id)

	if deleteErr != nil {
		return deleteErr
	}
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

func getIdFromParams(r *http.Request) (int64, error) {
	id := r.PathValue("id")

	parsedId, err := strconv.ParseInt(id, 10, 64)

	return parsedId, err

}
