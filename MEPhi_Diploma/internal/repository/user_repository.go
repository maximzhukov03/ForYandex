package repository

import (
    "database/sql"
    "errors"
    "golang/myapp/internal/models"
)

type UserRepository interface{
    Create(user *models.User) error
    GetByEmail(email string) (*models.User, error)
}

type SQLiteUserRepository struct{
    db *sql.DB
}

func NewSQLiteUserRepository(db *sql.DB) *SQLiteUserRepository{
    return &SQLiteUserRepository{
        db: db,
    }
}

func (r *SQLiteUserRepository) Create(user *models.User) error{
    query := `INSERT INTO users(email, password_hash, created_at) VALUES(?, ?, ?)`
    res, err := r.db.Exec(query, user.Email, user.PasswordHash, user.CreatedAt)
    if err != nil{
        return err
    }
    id, err := res.LastInsertId()
    if err != nil{
        return err
    }
    user.ID = id
    return nil
}

func (r *SQLiteUserRepository) GetByEmail(email string) (*models.User, error){
    query := `SELECT id, email, password_hash, created_at FROM users WHERE email = ?`
    row := r.db.QueryRow(query, email)
    var user models.User
    if err := row.Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt); err != nil{
        if errors.Is(err, sql.ErrNoRows){
            return nil, nil
        }
        return nil, err
    }
    return &user, nil
}