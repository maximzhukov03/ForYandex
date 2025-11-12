package repository

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func (bf *BrowseFiles) List(path string) ([]Object, error){
	allObjects, err := os.ReadDir(path)
	if err != nil{
		return nil, fmt.Errorf("try to read dir: %w", err)
	}
	
	objectSlice := make([]Object, 0, len(allObjects))
	for _, obj := range allObjects{
		info, err := obj.Info()
		if err != nil{
			return nil, fmt.Errorf("info object: %w", err)
		}
		if info.IsDir(){
			object := Object{
				Name: info.Name(),
				Size: 0,
				Type: TypeDir,
				TimeModification: info.ModTime(),
			}
			objectSlice = append(objectSlice, object)
		}else{
			object := Object{
				Name: info.Name(),
				Size: int(info.Size()),
				Type: TypeFile,
				TimeModification: info.ModTime(),
			}
			objectSlice = append(objectSlice, object)
		}
	}
	sort.SliceStable(objectSlice, func(i, j int) bool {
		if objectSlice[i].Type != objectSlice[j].Type{
			return objectSlice[i].Type == TypeDir
		}
		return strings.ToLower(objectSlice[i].Name) < strings.ToLower(objectSlice[j].Name)
	})
	return objectSlice, nil
}

func (bf *BrowseFiles) Exist(path string) (ex bool, Dir bool, err error){
	stat, err := os.Stat(path)
	if err != nil{
		if os.IsNotExist(err){
			return false, false, nil
		}
		return false, false, err
	}
	return true, stat.IsDir(), nil
}

func (bf *BrowseFiles) MakeDir(path, name string) error{
	p := filepath.Clean(filepath.Join(path, name))
	return os.Mkdir(p, accessRights)
}

func (bf *BrowseFiles) Rename(path, name, newName string) error{
	
	oldPath := filepath.Join(path, name)
	newPath := filepath.Join(path, newName)

	return os.Rename(oldPath, newPath)
}

func (bf *BrowseFiles) Delete(path, name string) error{
	fullPath := filepath.Clean(filepath.Join(path, name))
	err := os.RemoveAll(fullPath)
	if err != nil{
		return fmt.Errorf("remove error: %w", err)
	}
	return nil
}