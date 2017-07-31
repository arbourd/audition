package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type APIService struct {
	db Database
}

type CreateMessageRequest struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

func (a *APIService) ListMessages(w http.ResponseWriter, r *http.Request) {
	msgs, err := a.db.ListMessages()
	if err != nil {
		renderError(w, http.StatusInternalServerError, err)
		return
	}

	renderResponse(w, http.StatusOK, msgs)
}

func (a *APIService) GetMessage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		renderError(w, http.StatusInternalServerError, err)
		return
	}

	msg, err := a.db.GetMessage(int(id))
	if err != nil {
		if _, ok := err.(*NotFoundError); ok {
			renderError(w, http.StatusNotFound, err)
			return
		}

		renderError(w, http.StatusInternalServerError, err)
		return
	}
	renderResponse(w, http.StatusOK, msg)
}

func (a *APIService) CreateMessage(w http.ResponseWriter, r *http.Request) {
	var createMsg CreateMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&createMsg); err != nil {
		renderError(w, http.StatusInternalServerError, err)
		return
	}
	if createMsg.Message == "" {
		err := errors.New("'message' param cannot be an empty string")
		renderError(w, http.StatusBadRequest, err)
		return
	}

	m := &Message{
		Message: createMsg.Message,
	}
	msg, err := a.db.CreateMessage(m)
	if err != nil {
		renderError(w, http.StatusUnprocessableEntity, err)
	}

	renderResponse(w, http.StatusCreated, msg)
}

func (a *APIService) DeleteMessage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		renderError(w, http.StatusInternalServerError, err)
		return
	}

	err = a.db.DeleteMessage(int(id))
	if err != nil {
		if _, ok := err.(*NotFoundError); ok {
			renderError(w, http.StatusNotFound, err)
			return
		}

		renderError(w, http.StatusInternalServerError, err)
		return
	}
	renderResponse(w, http.StatusNoContent, nil)
}

func renderResponse(w http.ResponseWriter, code int, resp interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if resp != nil {
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("error encoding json: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func renderError(w http.ResponseWriter, code int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	resp := &ErrorResponse{
		Error:   http.StatusText(code),
		Message: err.Error(),
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("error encoding json: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Printf("%v: %v", http.StatusText(code), err.Error())
}
