package service

import (
    "errors"
    "time"

    "golang.org/x/crypto/bcrypt"

    "golang/myapp/internal/models"
    "golang/myapp/internal/repository"
)

type UserService interface{
    Register(email, password string) (*models.User, error)
    Authenticate(email, password string) (*models.User, error)
}

type userService struct{
    repo repository.UserRepository
    cost int
}

func NewUserService(repo repository.UserRepository, bcryptCost int) UserService{
    return &userService{repo: repo, cost: bcryptCost}
}

func (s *userService) Register(email, password string) (*models.User, error){
    exists, err := s.repo.GetByEmail(email)
    if err != nil{
        return nil, err
    }
    if exists != nil{
        return nil, errors.New("user already exists")
    }
    hash, err := bcrypt.GenerateFromPassword([]byte(password), s.cost)
    if err != nil{
        return nil, err
    }
    user := &models.User{
        Email:        email,
        PasswordHash: string(hash),
        CreatedAt:    time.Now(),
    }
    if err := s.repo.Create(user); err != nil{
        return nil, err
    }
    user.PasswordHash = ""
    return user, nil
}

func (s *userService) Authenticate(email, password string) (*models.User, error){
    user, err := s.repo.GetByEmail(email)
    if err != nil{
        return nil, err
    }
    if user == nil{
        return nil, errors.New("invalid credentials")
    }
    if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil{
        return nil, errors.New("invalid credentials")
    }
    user.PasswordHash = ""
    return user, nil
}