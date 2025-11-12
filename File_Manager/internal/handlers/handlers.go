package handlers

import (
	"browserfiles/test/internal/service"
	"encoding/json"
	"errors"
	"net/http"
	"os"
)

type Handler struct{
	service service.FileServicer
}

func NewHandler(serv service.FileServicer) *Handler{ 
	return &Handler{service: serv}
}

func writeJSON(w http.ResponseWriter, status int, v any){
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}


func mapErr(w http.ResponseWriter, err error){
	switch{
	case errors.Is(err, service.ErrOutRoot):
		writeJSON(w, http.StatusForbidden, map[string]string{"error": err.Error()})
	case errors.Is(err, service.ErrEmptyName),
		errors.Is(err, service.ErrInvalidName):
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
	case os.IsNotExist(err):
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
	default:
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
}
