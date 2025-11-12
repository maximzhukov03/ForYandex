package main

import (
	"log"
	"net/http"
	"os"
	repo "browserfiles/test/internal/repository"
	service "browserfiles/test/internal/service"
	handlers "browserfiles/test/internal/handlers"
)

func main() {
	rootDir := "./storage"
	if _, err := os.Stat(rootDir); os.IsNotExist(err) {
		if err := os.MkdirAll(rootDir, 0755); err != nil {
			log.Fatalf("cannot create root directory: %v", err)
		}
	}
	repo, err := repo.NewBrowseFiles(rootDir)
	if err != nil {
		log.Fatalf("init repository: %v", err)
	}
	serv := service.NewFileService(repo)
	h := handlers.NewHandler(serv)
	http.HandleFunc("/api/list", h.List)
	http.HandleFunc("/api/mkdir", h.Mkdir)
	http.HandleFunc("/api/rename", h.Rename)
	http.HandleFunc("/api/delete", h.Delete)
	http.HandleFunc("/api/upload", h.Upload)
	http.HandleFunc("/api/download", h.Download)
	fileServer := http.FileServer(http.Dir("./site"))
	http.Handle("/", fileServer)
	addr := ":8080"
	log.Printf("Server started on %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}