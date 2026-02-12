package repository

import (
	"context"
	"dailyPlanner/internal/database"
	"dailyPlanner/internal/models"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"slices"
)

type userRepository struct {
	db *database.DB
}

func NewUserRepository(db *database.DB) *userRepository {
	return &userRepository{db}
}

func (r userRepository) CreateUser(ctx context.Context, user *models.User, password string) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("password hashing error: %w", err)
	}

	user.UserId = uuid.New().String()
	user.PasswordHash = string(passwordHash)

	query := `
		INSERT INTO users (user_id, user_name, email, password_hash, role)
		VALUES (:user_id, :user_name, :email, :password_hash, :role)
	`

	_, err = r.db.NamedExecContext(ctx, query, user)
	if err != nil {
		return fmt.Errorf("error when creating a user when accessing the database: %w", err)
	}

	return nil
}

func (r userRepository) GetUserById(ctx context.Context, userId string) (*models.User, error) {
	query := `SELECT * FROM users WHERE user_id = $1`

	var user models.User
	err := r.db.GetContext(ctx, &user, query, userId)
	if err != nil {
		return nil, fmt.Errorf("error accessing the database when receiving a user: %w", err)
	}

	return &user, nil
}

func (r userRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `SELECT * FROM users WHERE email = $1`

	var user models.User
	err := r.db.GetContext(ctx, &user, query, email)
	if err != nil {
		return nil, fmt.Errorf("error accessing the database when receiving a user by email: %w", err)
	}

	return &user, nil
}

func (r userRepository) VerifyPassword(ctx context.Context, email string, password string) (*models.User, error) {
	user, err := r.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("error receiving the user when verifying the password: %w", err)
	}

	// checking that the password hash is the same
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("error when comparing password: %w", err)
	}

	return user, nil
}

func (r userRepository) UpdateUsername(ctx context.Context, email, newUserName, password string) error {
	user, err := r.VerifyPassword(ctx, email, password)
	if err != nil {
		return fmt.Errorf("error receiving the user when updating the user_name: %w", err)
	}

	user.UserName = newUserName

	query := `UPDATE users 
			 SET user_name = :user_name
             WHERE user_id = :user_id`

	result, err := r.db.NamedExecContext(ctx, query, user)
	if err != nil {
		return fmt.Errorf("error when updating the user_name: %w", err)
	}

	if !CheckUpdate(result) {
		return fmt.Errorf("error when updating the user_name: no rows were affected")
	}

	return nil
}

func (r userRepository) UpdatePassword(ctx context.Context, email, password, newPassword string) error {
	user, err := r.VerifyPassword(ctx, email, password)
	if err != nil {
		return fmt.Errorf("error receiving the user when updating the password: %w", err)
	}

	newPasswordHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("error when updating the password: %w", err)
	}

	user.PasswordHash = string(newPasswordHash)

	query := `UPDATE users 
			 SET password_hash = :password_hash
             WHERE user_id = :user_id`

	result, err := r.db.NamedExecContext(ctx, query, user)
	if err != nil {
		return fmt.Errorf("error when updating the password: %w", err)
	}

	if !CheckUpdate(result) {
		return fmt.Errorf("error when updating the password: no rows were affected")
	}

	return nil
}

func (r userRepository) AppointmentModerator(ctx context.Context, email, role, password string) error {
	user, err := r.VerifyPassword(ctx, email, password)
	if err != nil {
		return fmt.Errorf("error receiving the user when appointment moderator: %w", err)
	}

	acceptableRoles := []string{"User", "Admin"}

	if !slices.Contains(acceptableRoles, role) {
		return fmt.Errorf("replacing with a non-existing role")
	}

	user.Role = role

	query := `UPDATE users 
			 SET role = :role
             WHERE user_id = :user_id`

	result, err := r.db.NamedExecContext(ctx, query, user)
	if err != nil {
		return fmt.Errorf("error when updating the role: %w", err)
	}

	if !CheckUpdate(result) {
		return fmt.Errorf("error when updating the role: no rows were affected")
	}
	return nil
}

func (r userRepository) DeleteUser(ctx context.Context, userId string) error {
	query := `DELETE FROM users WHERE user_id = $1`

	result, err := r.db.ExecContext(ctx, query, userId)
	if err != nil {
		return fmt.Errorf("error when deleting the user: %w", err)
	}

	if !CheckUpdate(result) {
		return fmt.Errorf("error when deleting the user: no rows were affected")
	}
	return nil
}

func CheckUpdate(result sql.Result) bool {
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false
	}

	if rowsAffected == 0 {
		return false
	}
	return true
}
