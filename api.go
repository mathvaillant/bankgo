package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type APIServer struct {
	listenAddr string
	storage    Storage
}

func NewApiServer(listenAddr string, store Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		storage:    store,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/accounts", makeHTTPHandleFunc(s.handleGetAccounts)).Methods("GET")
	router.HandleFunc("/accounts", makeHTTPHandleFunc(s.handleCreateAccount)).Methods("POST")
	router.HandleFunc("/accounts", makeHTTPHandleFunc(s.handleDeleteAccount)).Methods("DELETE")

	router.HandleFunc("/accounts/{id}", makeHTTPHandleFunc(s.handleGetAccountById)).Methods("GET")
	router.HandleFunc("/accounts/{id}", makeHTTPHandleFunc(s.handleDeleteAccount)).Methods("DELETE")

	router.HandleFunc("/transfer", makeHTTPHandleFunc(s.handleTransfer)).Methods("POST")

	log.Println("JSON api running on port: ", s.listenAddr)

	http.ListenAndServe(s.listenAddr, router)
}

func (s *APIServer) handleGetAccounts(w http.ResponseWriter, r *http.Request) error {
	accounts, err := s.storage.GetAccounts()

	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, accounts)
}

func (s *APIServer) handleGetAccountById(w http.ResponseWriter, r *http.Request) error {
	id, err := getIDFromParams(r)

	if err != nil {
		return err
	}

	account, err := s.storage.GetAccountByID(id)

	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, account)
}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	id, err := getIDFromParams(r)

	if err != nil {
		return err
	}

	if err := s.storage.DeleteAccount(id); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusNoContent, map[string]int64{"deleted": id})
}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	createAccountReq := &CreateAccountRequest{}

	if err := json.NewDecoder(r.Body).Decode(createAccountReq); err != nil {
		return err
	}

	account := NewAccount(createAccountReq.FirstName, createAccountReq.LastName)

	if err := s.storage.CreateAccount(account); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, account)
}

func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func WriteJSON(w http.ResponseWriter, statusCode int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(v)
}

type ApiError struct {
	Error string `json:"error"`
}
type ApiFunc func(http.ResponseWriter, *http.Request) error

func makeHTTPHandleFunc(f ApiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

func getIDFromParams(r *http.Request) (int64, error) {
	idStr := mux.Vars(r)["id"]
	idInt, err := strconv.Atoi(idStr)
	id := int64(idInt)

	if err != nil {
		return id, fmt.Errorf("invalid id given %s", idStr)
	}
	return id, nil
}
