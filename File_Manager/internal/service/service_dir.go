package service

import (
	"browserfiles/test/internal/repository"
	"strings"
)

func (s *FileService) Exist(logicPath string) (bool, bool, error){
	jS, err := s.joinSafe(logicPath)
	if err != nil{
		return  false, false, err
	}

	return  s.repo.Exist(jS)
}

func (s *FileService) List(logicalPath string) ([]repository.Object, error){
	p, err := s.joinSafe(logicalPath)
	if err != nil{
		return nil, err
	}
	return s.repo.List(p)
}

func (s *FileService) MakeDir(logicalPath, name string) error{
	name = strings.TrimSpace(name)
	if name == ""{
		return ErrEmptyName
	}
	p, err := s.joinSafe(logicalPath)
	if err != nil{
		return err
	}
	return s.repo.MakeDir(p, name)
}

func (s *FileService) Rename(logicalPath, name, newName string) error{
	name = strings.TrimSpace(name)
	newName = strings.TrimSpace(newName)
	if name == "" || newName == "" || newName == name{
		return ErrInvalidName
	}
	p, err := s.joinSafe(logicalPath)
	if err != nil{
		return err
	}
	return s.repo.Rename(p, name, newName)
}

func (s *FileService) Delete(logicalPath, name string) error{
	name = strings.TrimSpace(name)
	if name == ""{
		return ErrEmptyName
	}
	p, err := s.joinSafe(logicalPath)
	if err != nil{
		return err
	}
	return s.repo.Delete(p, name)
}