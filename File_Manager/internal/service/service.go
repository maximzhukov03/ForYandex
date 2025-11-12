package service

import (
	repository "browserfiles/test/internal/repository"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var (
	ErrOutRoot = errors.New("path outside of root")
	ErrEmptyName = errors.New("empty name")
	ErrInvalidName = errors.New("invalid new name")
)

type FileService struct{
	repo repository.Repository
	root string
}

type FileServicer interface{
	List(path string) ([]repository.Object, error)
	MakeDir(path, name string) error
	Rename(path, name, newName string) error
	Delete(path, name string) error
	Save(path, name string, r io.Reader) (int64, error)
	Download(path string) (*os.File, string, string, error)
	Exist(logicPath string) (ex bool, isDir bool, err error)
}

func NewFileService(repo repository.Repository) *FileService{
	return &FileService{
		repo: repo,
		root: repo.RootPath(),
	}
}

func (s *FileService) joinSafe(logicalPath string, parts ...string) (string, error){
	lp := strings.TrimSpace(logicalPath)
	if lp == ""{
		lp = "."
	}
	joined := filepath.Join(append([]string{s.root, lp}, parts...)...)
	clean := filepath.Clean(joined)

	if !strings.HasPrefix(clean+string(filepath.Separator), s.root+string(filepath.Separator)){
		return "", ErrOutRoot
	}
	return clean, nil
}



