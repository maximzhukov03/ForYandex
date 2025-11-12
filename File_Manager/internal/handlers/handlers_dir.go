package handlers

import (
	"encoding/json"
	"net/http"
)

func (h *Handler) List(w http.ResponseWriter, r *http.Request){
	path := r.URL.Query().Get("path")
	objs, err := h.service.List(path)
	if err != nil{
		mapErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"items": objs,
	})
}

func (h *Handler) Mkdir(w http.ResponseWriter, r *http.Request){
	var req dirReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil{
		writeJSON(w, http.StatusBadRequest, map[string]string{"error":"invalid json"})
		return
	}
	if err := h.service.MakeDir(req.Path, req.Name); err != nil{
		mapErr(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]string{"status":"ok"})
}

func (h *Handler) Rename(w http.ResponseWriter, r *http.Request){
	var req renameReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil{
		writeJSON(w, http.StatusBadRequest, map[string]string{"error":"invalid json"})
		return
	}
	if err := h.service.Rename(req.Path, req.Name, req.NewName); err != nil{
		mapErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status":"ok"})
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request){
	var req deleteReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil{
		writeJSON(w, http.StatusBadRequest, map[string]string{"error":"invalid json"})
		return
	}
	if err := h.service.Delete(req.Path, req.Name); err != nil{
		mapErr(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status":"ok"})
}