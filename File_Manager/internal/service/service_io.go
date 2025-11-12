package service

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func (s *FileService) Save(logicalPath, name string, r io.Reader) (int64, error){
	name = strings.TrimSpace(name)
	if name == ""{
		return 0, ErrEmptyName
	}
	p, err := s.joinSafe(logicalPath)
	if err != nil{
		return 0, err
	}
	return s.repo.Save(p, name, r)
}

func (s *FileService) Download(logicalPath string) (*os.File, string, string, error){
	p, err := s.joinSafe(logicalPath)
	if err != nil{
		return nil, "", "", err
	}
	ex, isDir, err := s.repo.Exist(p)
	if err != nil{
		return nil, "", "", err
	}
	if !ex{
		return nil, "", "", os.ErrNotExist
	}
	if isDir{
		return nil, "", "", errors.New("cannot download directories")
	}
	f, mime, err := s.repo.Download(p)
	if err != nil{
		return nil, "", "", err
	}
	filename := filepath.Base(p)
	return f, mime, filename, nil
}