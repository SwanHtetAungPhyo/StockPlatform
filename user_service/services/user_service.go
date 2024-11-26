package service

import (
	"context"
	"errors"
	"fmt"
	"time"
	"github.com/SwanHtetAungPhyo/swan_lib"
	"github.com/SwanHtetAungPhyo/user-service/models"
	"github.com/SwanHtetAungPhyo/user-service/repository"
	"golang.org/x/crypto/bcrypt"
)

var repo repository.UserRepository
var secret_key string 
var time_frame time.Duration
func SetRepository(r repository.UserRepository,key string, expire time.Duration ) {
	repo = r
	secret_key = key
	time_frame = expire
}

func SignUp(user *models.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("error hashing password: %v", err)
	}
	user.Password = string(hashedPassword)

	err = repo.SaveUser(*user)
	if err != nil {
		return fmt.Errorf("error saving user: %v", err)
	}

	return nil
}

func SignIn(user *models.User, timeoutCtx context.Context) (string, error) {

	storedUser, err := repo.FindUserByEmail(user.Email) 
	if err != nil {
		return "", fmt.Errorf("error finding user: %v", err)
	}


	if storedUser == nil {
		return "", errors.New("user not found")
	}


	err = bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(user.Password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}
	jwtManager := swan_lib.NewJWTManager(secret_key, time_frame)


	customClaims := map[string]any{
		"role": "admin", }
	authToken, err := jwtManager.GenerateToken(storedUser.Name, customClaims)
	if err != nil {
		return "", fmt.Errorf("error generating auth token: %v", err)
	}

	return authToken, nil
}

func Deposit(userEmail string, amount float64, ctx context.Context) (float64, error) {
	user, err := repo.FindUserByEmail(userEmail)
	if err != nil {
		return 0, fmt.Errorf("failed to find user by email: %w", err)
	}
	if user == nil {
		return 0, fmt.Errorf("user not found")
	}


	var balance *models.Balance
	balance,err = repo.FindBalanceByUserID(user.ID)
	if err != nil {
		return 0, fmt.Errorf("failed to find balance: %w", err)
	}

	newBalance := balance.Amount + amount
	err = repo.UpdateBalance(user.ID, newBalance)  
	if err != nil {
		return 0, fmt.Errorf("failed to update balance: %w", err)
	}

	return newBalance, nil
}

func GetBalance(userEmail string, timeout context.Context) (float64, error) {
	user, err := repo.FindUserByEmail( userEmail)
	if err != nil {
		return 0, fmt.Errorf("error finding user: %w", err)
	}

	var balance *models.Balance
	balance,err = repo.FindBalanceByUserID(user.ID)
	if err != nil {
		return 0, fmt.Errorf("failed to find balance: %w", err)
	}
	return balance.Amount, nil
}