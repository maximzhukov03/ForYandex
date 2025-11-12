package repository

import (
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
)

func (bf *BrowseFiles) Save(path, name string, r io.Reader) (int64, error){
	distPath := filepath.Join(path, name)
	if bf.Opt.ConflictingName == true{
		return 0, fmt.Errorf("file or dir name '%s' not unique", name)
	}
	tempFile, err := os.CreateTemp(path, ".tempFile-*")
	if err != nil{
		return 0, fmt.Errorf("create temp file: %w", err)
	}
	written, err := io.Copy(tempFile, r)
	if err != nil{
		err := tempFile.Close()
		if err != nil{
			return written, fmt.Errorf("close temporary file: %w", err)
		}
		return written, err 
	}
	err = tempFile.Sync()
	if err != nil{
		err := tempFile.Close()
		if err != nil{
			return written, fmt.Errorf("error close temporary file")
		}
		return written, fmt.Errorf("error sync commits the current contents: %w", err) 
	}
	err = tempFile.Close()
	if err != nil{
		return written, fmt.Errorf("error of close temp file: %w", err)
	}

	if bf.Opt.ConflictingName{
		err := os.Remove(distPath)
		if err != nil{
			return written, fmt.Errorf("error of deleting temporary file")
		}
	}

	_, err = os.Stat(distPath)
	if err == nil{
		_ = os.Remove(distPath)
	}

	err = os.Rename(tempFile.Name(), distPath)
	if err != nil{
		_ = os.Remove(tempFile.Name())
		return written, fmt.Errorf("rename temp file: %w", err)
	}
	return written, nil
}

func (bf *BrowseFiles) Download(path string) (*os.File, string, error){
	stat, err := os.Stat(path)
	if err != nil{
		return nil, "" ,fmt.Errorf("stat info error: %w",err)
	}
	if stat.IsDir(){
		return nil, "", errors.New("error download, path is dir")
	}
	op, err := os.Open(path)
	if err != nil{
		return nil, "", fmt.Errorf("error open os path: %w", err)
	}
	mi := mime.TypeByExtension(filepath.Ext(path))
	if mi == ""{
		buffer := make([]byte, 512)
		n, err := op.Read(buffer)
		if err != nil && err != io.EOF{
			op.Close()
			return nil, "", fmt.Errorf("read file for mime: %w", err)
		}
		_,_ = op.Seek(0, io.SeekStart)
		mi = http.DetectContentType(buffer[:n])
	}
	if mi == ""{
		mi = "application/octet-stream"
	}
	return op, mi, nil
}