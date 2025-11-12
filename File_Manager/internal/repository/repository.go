package repository

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type Repository interface{
	List(path string) ([]Object, error)
	MakeDir(path, name string) error
	Rename(path, name, newName string) error
	Delete(path, name string) error
	Save(path, name string, r io.Reader) (int64, error)
	Download(path string) (*os.File, string, error)
	Exist(path string) (ex bool, Dir bool, err error)
	RootPath() string
}


func NewBrowseFiles(root string) (*BrowseFiles, error){
	if root == ""{
		return nil, errors.New("empty root")
	}

	abs, err := filepath.Abs(root)
	if err != nil{
		return nil, fmt.Errorf("abs filepath error: %w", err)
	}

	stat, err := os.Stat(abs)
	if err != nil{
		return nil, fmt.Errorf("stat info error: %w", err)
	}
	if !stat.IsDir(){
		return nil, fmt.Errorf("root isnt dir: %s", abs)
	}
	return &BrowseFiles{
		Root: filepath.Clean(abs),
	}, nil

}

func (bf *BrowseFiles) RootPath() string{
	return bf.Root
}
