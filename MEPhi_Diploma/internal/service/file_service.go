package service

import (
	"fmt"
	"io"
	"time"

	"golang/myapp/internal/models"
	"golang/myapp/internal/repository"

	"github.com/google/uuid"
)

type FileService interface{
    Upload(userID int64, name string, reader io.Reader, size int64, contentType string) (*models.File, error)
    List(userID int64) ([]*models.File, error)
}

type fileService struct{
    repo   repository.FileRepository
    store  repository.StorageClient
    bucket string
    urlTTL time.Duration
}

func NewFileService(repo repository.FileRepository, store repository.StorageClient, bucket string, urlTTL time.Duration) FileService{
    return &fileService{
        repo: repo, 
        store: store, 
        bucket: bucket, 
        urlTTL: urlTTL,
    }
}

func (s *fileService) Upload(
    userID int64, 
    name string, 
    reader io.Reader, 
    size int64, 
    contentType string) (*models.File, error){
    objectName := fmt.Sprintf("%s_%s", uuid.New().String(), name)
    if err := s.store.Upload(s.bucket, objectName, reader, size, contentType); err != nil{
        return nil, err
    }

    file := &models.File{
        UserID:     userID,
        Name:       name,
        Size:       size,
        Bucket:     s.bucket,
        ObjectName: objectName,
        UploadedAt: time.Now(),
    }
    if err := s.repo.Save(file); err != nil{
        return nil, err
    }
    url, err := s.store.PresignedURL(s.bucket, objectName, s.urlTTL)
    if err != nil{
        return nil, err
    }
    file.URL = url
    return file, nil
}

func (s *fileService) List(userID int64) ([]*models.File, error){
    files, err := s.repo.ListByUserID(userID)
    if err != nil{
        return nil, err
    }
    for _, f := range files{
        url, err := s.store.PresignedURL(f.Bucket, f.ObjectName, s.urlTTL)
        if err != nil{
            return nil, err
        }
        f.URL = url
    }
    return files, nil
}