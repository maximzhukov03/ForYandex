package service

import (
	proto "content/api/gen/content/v1"
	"content/api/internal/storage"
	"context"
	"fmt"
	"io"
)

type FileService struct{
	proto.UnimplementedContentServiceServer
	storage storage.FileStorage
}

func NewFileService(storage storage.FileStorage) *FileService{
	return  &FileService{
		storage: storage,
	}
}

type FileServicer interface{
	Upload(flow proto.ContentService_UploadServer) error
	Download(req *proto.DownloadRequest, flow proto.ContentService_DownloadServer) error
	List(ctx context.Context, req *proto.UploadRequest) (*proto.ListResponse, error)
}

func (fs *FileService) Upload(flow proto.ContentService_UploadServer) error{

	ctx := flow.Context()
	reader, writer := io.Pipe()
	doneCh := make(chan struct{})

	part, err := flow.Recv()
	if err != nil{
		return  fmt.Errorf("Error Receive Part: %w", err)
	}
	
	fileName := part.FileName

	var err1 error
	var	file storage.File

	go func(){
		defer close(doneCh)
		file, err1 = fs.storage.Upload(ctx, fileName, reader) 
	}()

	if len(part.Data) > 0{
		if _, err := writer.Write(part.Data); err != nil {
			writer.CloseWithError(err)
			return fmt.Errorf("Error Write Part: %w", err)
		}
	}

	for {
		nextPart, err := flow.Recv()
		if err == io.EOF{
			writer.Close()
			break
		}
		if err != nil {
			writer.CloseWithError(err)
			return fmt.Errorf("Error Receive Part: %w", err)
		}
		if len(nextPart.Data) > 0{
			if _, err := writer.Write(nextPart.Data); err != nil{
				writer.CloseWithError(err)
				return fmt.Errorf("Error Write Part: %w", err)
			}
		}
	}
	<-doneCh

	if err1 != nil{
		return fmt.Errorf("Error Save File: %w", err1)
	}

	resp := &proto.UploadResponse{
		Id:         file.ID,
		FileName:   file.File_name,
		Size:       file.Size,
		CreatedAt:  file.Created.Unix(),
		ControlSum: file.ControlSum,
	}
	return flow.SendAndClose(resp)
}

func (fs *FileService) Download(req *proto.DownloadRequest, flow proto.ContentService_DownloadServer) error{
	ctx := flow.Context()
	_, reader, err := fs.storage.Open(ctx, req.Id)
	if err != nil{
		return fmt.Errorf("Error Open File: %w", err)
	}
	defer reader.Close()
	
	buffer := make([]byte, 1024*64)
	for{
		n, err := reader.Read(buffer)
		if err == io.EOF{
			break
		}
		if err != nil{
			return  fmt.Errorf("Error Read File: %w", err)
		}
		err = flow.Send(&proto.DownloadResponse{Data: buffer[:n]})
		if err != nil{
			return fmt.Errorf("Error Send: %w", err)
		}
	}
	return nil
}

func (fs *FileService) List(ctx context.Context, req *proto.ListRequest) (*proto.ListResponse, error){
	files, nextIndx, err := fs.storage.List(ctx, int(req.Count), req.Indx)
	if err != nil{
		return  nil, fmt.Errorf("Error From List: %w", err)
	}

	resp := &proto.ListResponse{
		RightIndx: nextIndx,
	}

	for _, file := range files{
		timeCr := file.Created.Unix()
		timeUp := file.Updated.Unix() 
		resp.Item = append(resp.Item, &proto.File{
			Id: file.ID,
			FileName: file.File_name,
			Size: file.Size,
			Created: timeCr,
			Updated: timeUp,
			ControlSum: file.ControlSum,
		})
	}

	return  resp, nil
}