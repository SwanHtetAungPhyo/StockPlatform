package repository

import (
	"fmt"

	"github.com/SwanHtetAungPhyo/user-service/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	SaveUser(user models.User) error
	FindUserByEmail( email string) (*models.User, error)
	UpdateBalance(userID uint, newBalance float64) error
	FindBalanceByUserID(userID uint) (*models.Balance, error)
}

type userRepository struct {
	DB *gorm.DB
}

func NewRepository(db *gorm.DB) UserRepository {
	return &userRepository{DB: db}
}

func (repo *userRepository) SaveUser(user models.User) error {
	balance := new(models.Balance)
	tx := repo.DB.Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to start transaction: %w", tx.Error)
	}

	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback() 
		return fmt.Errorf("failed to save user: %w", err)
	}


	balance.UserID = user.ID


	if err := tx.Create(&balance).Error; err != nil {
		tx.Rollback() 
		return fmt.Errorf("failed to save balance: %w", err)
	}


	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}


func (repo *userRepository) FindUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := repo.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find user by email: %w", err)
	}
	return &user, nil
}


func (repo *userRepository) UpdateBalance(userID uint, newBalance float64) error {
	balance, err := repo.FindBalanceByUserID(userID)
	if err != nil {
		return fmt.Errorf("failed to find balance: %w", err)
	}
	if err := repo.UpdateBalanceAmount(balance, newBalance); err != nil {
		return fmt.Errorf("failed to update balance: %w", err)
	}

	return nil
}
func (repo *userRepository) UpdateBalanceAmount(balance *models.Balance, newBalance float64) error {
	balance.Amount = newBalance
	if err := repo.DB.Save(balance).Error; err != nil {
		return fmt.Errorf("failed to update balance: %w", err)
	}

	return nil
}

func (repo *userRepository) FindBalanceByUserID(userID uint) (*models.Balance, error) {
	var balance models.Balance
	err := repo.DB.Where("user_id = ?", userID).First(&balance).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("balance not found for user %d", userID)
		}
		return nil, fmt.Errorf("failed to find balance: %w", err)
	}

	return &balance, nil
}
