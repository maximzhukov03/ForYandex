package storage

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
	_ "modernc.org/sqlite"
	"github.com/oklog/ulid/v2"
)

type File struct {
	ID          string
	File_name   string
	Size        int64
	Path 		string
	ControlSum string
	Updated     time.Time
	Created     time.Time
}

type FileStorageSQlite struct {
	db  *sql.DB
	dir string
}

func NewFileStorageSQlite(path, dir string) (*FileStorageSQlite, error) {
	d, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("Error SQLite open process: %w", err)
	}

	storage := &FileStorageSQlite{
		db:  d,
		dir: dir,
	}
	err = storage.createTable()
	if err != nil {
		return nil, fmt.Errorf("Error Create Table: %w", err)
	}
	return storage, nil

}

func (f *FileStorageSQlite) createTable() error {
	_, err := f.db.Exec(
		`CREATE TABLE IF NOT EXISTS files 
		(id TEXT PRIMARY KEY, 
		file_name TEXT NOT NULL, 
		path TEXT NOT NULL, 
		size INTEGER NOT NULL, 
		control_sum TEXT NOT NULL, 
		created INTEGER NOT NULL, 
		updated INTEGER NOT NULL); 
		CREATE INDEX IF NOT EXISTS idx_files_created ON files(created DESC);`)
	return err
}

type FileStorage interface {
	Upload(ctx context.Context, filename string, reader io.Reader) (File, error)
	Open(ctx context.Context, fileId string) (File, io.ReadCloser, error)
	List(ctx context.Context, count int, indx string) ([]File, string, error)
}

func CreatePath(dir, id string) (tmpPath, resPath string, err error){
	directory := filepath.Join(dir, "imgs")
	err = os.MkdirAll(directory, 0755)
	if err != nil{
		return "", "", fmt.Errorf("Error Create Path(MkdirAll): %w", err)
	}
	tmp := fmt.Sprintf("%s.tmp", id)
	tmpPath = filepath.Join(directory, tmp)
	resPath = filepath.Join(directory, id)
	return  tmpPath, resPath, err
}

func (f *FileStorageSQlite) Upload(ctx context.Context, filename string, reader io.Reader) (File, error){
	t := time.Now()
	en := ulid.Monotonic(rand.Reader, 0)
	id := ulid.MustNew(ulid.Timestamp(t), en).String()
	tmpPath, resPath, err := CreatePath(f.dir, id)
	if err != nil{
		return File{}, err
	}
	res, err := os.Create(tmpPath)
	if err != nil{
		return File{}, fmt.Errorf("Error Create File(os.Create): %w", err)
	}
	defer res.Close()

	hash := sha256.New()
	size, err := io.Copy(io.MultiWriter(res, hash), reader)
	if err != nil{
		return File{}, fmt.Errorf("Error Get Size (io.Copy): %w", err)
	
	}
	control_sum := hex.EncodeToString(hash.Sum(nil))

	fi := File{
		ID: id,
		File_name: filename,
		Size: size,
		Path: resPath,
		ControlSum: control_sum,
		Updated: t,
		Created: t,

	}

	err = res.Sync()
	if err != nil{
		return  File{}, fmt.Errorf("Error Sync: %w", err)
	}

	err = os.Rename(tmpPath, resPath)
	if err != nil{
		return File{}, fmt.Errorf("Error Rename: %w", err)
	}

	transaction, err := f.db.BeginTx(ctx, nil)
	if err != nil{
		return fi, fmt.Errorf("Error Transaction(BegiTx): %w", err)
	}

	_, err = transaction.ExecContext(ctx, `INSERT INTO files(
	id, 
	file_name, 
	path, size,
	control_sum, 
	created, 
	update) 
	VALUE (?, ?, ?, ?, ?, ?, ?)`, 
	fi.ID, fi.File_name, fi.Path, fi.Size, fi.ControlSum, fi.Created.Unix(), fi.Updated.Unix())
	if err != nil{
		transaction.Rollback()
		os.Remove(resPath)
		return  fi, fmt.Errorf("Error Exect: %w", err)
	}

	transaction.Commit()
	return fi, nil
}

func (f *FileStorageSQlite)Open(ctx context.Context, fileId string) (File, io.ReadCloser, error){
	fi := File{}
	var created, updated int64
	fi.ID = fileId
	query := f.db.QueryRowContext(ctx, `SELECT file_name, path, size, control_sum, created, updated FROM files WHERE id = ?`, fileId)

	
	err := query.Scan(&fi.File_name, &fi.Path, &fi.Size, &fi.ControlSum, &created, &updated)
	if err != nil{
		return  fi, nil, fmt.Errorf("Error Query Scan: %w", err)
	}
	fi.Created = time.Unix(created, 0)
	fi.Updated = time.Unix(updated, 0)
	
	op, err := os.Open(fi.Path)
	if err != nil{
		return fi, nil, fmt.Errorf("Error os.Open file: %w", err)
	}
	return fi, op, nil 

}

func (f *FileStorageSQlite) List(ctx context.Context, count int, indx string) ([]File, string, error){
	var rows *sql.Rows
	var err error

	if indx == ""{
		rows, err = f.db.QueryContext(ctx, `SELECT id, file_name, path, size, control_sum, created, updated
		FROM files ORDER BY created DESC LIMIT ?`, count)
	} else {
		ti, err2 := time.Parse(time.RFC3339, indx)
		if err2 != nil{
			return nil, "", fmt.Errorf("Error indx: %w", err2)
		}
		rows, err = f.db.QueryContext(ctx, `SELECT id, file_name, path, size, control_sum, created, updated
		FROM files WHERE created < ? ORDER BY created DESC LIMIT ?`, ti.Unix(), count)
	}

	if err != nil{
		return nil, "", fmt.Errorf("Error Query: %w", err)
	}
	defer rows.Close()

	files := []File{}
	for rows.Next(){
		fi := File{}
		if err := rows.Scan(&fi.ID, &fi.File_name, &fi.Path, &fi.Size, &fi.ControlSum, &fi.Created, &fi.Updated); err != nil {
		return nil, "", err
		}
		files = append(files, fi)
	}

	nextIndx := ""
	if len(files) > 0{
		nextIndx = files[len(files)-1].Created.Format(time.RFC3339)
	}

	return files, nextIndx, nil
}