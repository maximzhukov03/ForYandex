package service

import (
    "errors"

    "golang/myapp/internal/models"
    "golang/myapp/internal/repository"
)

type AdminService interface{
    ListUsers() ([]*models.User, error)
    GetUserByID(id int64) (*models.User, error)
    UpdateUser(id int64, email, role string) (*models.User, error)
    DeleteUser(id int64) (bool, error)
    PromoteUser(id int64) (*models.User, error)
}

type adminService struct{
    repo repository.AdminRepository
}

func NewAdminService(repo repository.AdminRepository) AdminService{
    return &adminService{
        repo: repo,
    }
}

func (s *adminService) ListUsers() ([]*models.User, error){
    return s.repo.List()
}

func (s *adminService) GetUserByID(id int64) (*models.User, error){
    return s.repo.GetByID(id)
}

func (s *adminService) UpdateUser(id int64, email, role string) (*models.User, error){
    user, err := s.repo.GetByID(id)
    if err != nil{
        return nil, err
    }
    if user == nil{
        return nil, nil
    }
    user.Email = email
    user.Role = role
    if err := s.repo.Update(user); err != nil{
        return nil, err
    }
    return user, nil
}

func (s *adminService) DeleteUser(id int64) (bool, error){
    rows, err := s.repo.Delete(id)
    if err != nil{
        return false, err
    }
    return rows > 0, nil
}

func (s *adminService) PromoteUser(id int64) (*models.User, error){
    user, err := s.repo.GetByID(id)
    if err != nil{
        return nil, err
    }
    if user == nil{
        return nil, nil
    }
    if user.Role == "admin"{
        return user, errors.New("user is already admin")
    }
    user.Role = "admin"
    if err := s.repo.Update(user); err != nil{
        return nil, err
    }
    return user, nil
}