package models

import (
    "time"
)

type User struct{
    ID           int64     `db:"id" json:"id"`
    Email        string    `db:"email" json:"email"`
    PasswordHash string    `db:"password_hash" json:"-"`
    CreatedAt    time.Time `db:"created_at" json:"created_at"`
    Role         string    `db:"role" json:"role"`
}

type File struct{
    ID         int64     `db:"id" json:"id"`
    UserID     int64     `db:"user_id" json:"user_id"`
    Name       string    `db:"name" json:"name"`
    Size       int64     `db:"size" json:"size"`
    Bucket     string    `db:"bucket" json:"bucket"`
    ObjectName string    `db:"object_name" json:"object_name"`
    UploadedAt time.Time `db:"uploaded_at" json:"uploaded_at"`
    URL        string    `db:"-" json:"url,omitempty"`
}
