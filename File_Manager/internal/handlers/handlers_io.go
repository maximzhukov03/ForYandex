package handlers

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

func (h *Handler) Upload(w http.ResponseWriter, r *http.Request){
	path := r.URL.Query().Get("path")

	if err := r.ParseMultipartForm(32 << 20); err != nil{
		writeJSON(w, http.StatusBadRequest, map[string]string{"error":"multipart parse error"})
		return
	}
	file, header, err := r.FormFile("file")
	if err != nil{
		writeJSON(w, http.StatusBadRequest, map[string]string{"error":"missing file"})
		return
	}
	defer file.Close()

	name := r.FormValue("name")
	if strings.TrimSpace(name) == ""{
		name = header.Filename
	}

	if _, err := h.service.Save(path, name, file); err != nil{
		mapErr(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]string{"status":"ok"})
}

func (h *Handler) Download(w http.ResponseWriter, r *http.Request){
	path := r.URL.Query().Get("path")
	f, mime, filename, err := h.service.Download(path)
	if err != nil {
		mapErr(w, err)
		return
	}
	defer f.Close()

	w.Header().Set("Content-Type", mime)
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	_, _ = io.Copy(w, f)
}