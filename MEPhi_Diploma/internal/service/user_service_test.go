package service_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"golang/myapp/internal/models"
	"golang/myapp/internal/service"
)

type MockUserRepo struct { 
	mock.Mock 
}

func (m *MockUserRepo) Create(u *models.User) error{
	args := m.Called(u)
	return args.Error(0)
}
func (m *MockUserRepo) GetByEmail(email string) (*models.User, error){
	args := m.Called(email)
	user, _ := args.Get(0).(*models.User)
	return user, args.Error(1)
}
func (m *MockUserRepo) GetByID(id int64) (*models.User, error){ 
	return nil, nil 
}
func (m *MockUserRepo) Update(u *models.User) error{
	return nil
}
func (m *MockUserRepo) Delete(id int64) (int64, error){
	return 0, nil
}

const testBcryptCost = 12

func TestRegister_Success(t *testing.T){
	repo := new(MockUserRepo)
	us := service.NewUserService(repo, testBcryptCost)
	repo.On("GetByEmail", "maxim@exmail.com").Return(nil, nil)
	repo.On("Create", mock.Anything).Return(nil)
	user, err := us.Register("maxim@exmail.com", "Password123")
	assert.NoError(t, err)
	assert.Equal(t, "maxim@exmail.com", user.Email)
	repo.AssertExpectations(t)
}

func TestRegister_DuplicateEmail(t *testing.T){
	repo := new(MockUserRepo)
	us := service.NewUserService(repo, testBcryptCost)
	repo.On("GetByEmail", "max@exmail.com").Return(&models.User{}, nil)
	_, err := us.Register("max@exmail.com", "Password123")
	assert.EqualError(t, err, errors.New("user already exists").Error())
	repo.AssertCalled(t, "GetByEmail", "max@exmail.com")
}