package services

import (
	"fmt"

	"github.com/harshgupta9473/Social/pkg/models"
	"github.com/harshgupta9473/Social/pkg/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	UserRepo *repository.UserRepository
}

func NewAuthService(userRepo *repository.UserRepository) *AuthService {
	return &AuthService{UserRepo: userRepo}
}

func (auth *AuthService) Register(user models.TempUser) error {
	err := auth.UserRepo.CreateUserAccount(user)
	if err != nil {
		return err
	}
	return err
}

func ComparePassword(hashPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
}

func (auth *AuthService) SignIn(email, password string) (*models.User, error) {
	user, err := auth.UserRepo.GetUserByEmail(email)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	err = ComparePassword(user.Encrypted_Password, password)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("wrong email or password")
	}
	return user, nil
}
