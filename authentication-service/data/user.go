package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const dbTimeout = 3 * time.Second

type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name,omitempty"`
	LastName  string    `json:"last_name,omitempty"`
	Password  string    `json:"-"`
	Active    int       `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserModel will handle the db interactions making it modular and easy to plug in any
// other db orm like gorm
type UserModel struct {
	db *sql.DB
}

// NewUserModel injects the db object to the UserModel struct and inits it
func NewUserModel(dbPool *sql.DB) *UserModel {
	return &UserModel{
		db: dbPool,
	}
}

// func Insert() -> inserts a new user and returns the inserted ID
func (u *UserModel) Insert(user *User) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		return 0, fmt.Errorf("error hashing password %w", err)
	}

	stmt := `insert into users(email,first_name,last_name,password, user_active,created_at,updated_at)
	values($1,$2,$3,$4,$5,$6,$7) returning id`

	var insertedID int

	err = u.db.QueryRowContext(ctx, stmt,
		user.Email,
		user.FirstName,
		user.LastName,
		hashedPassword,
		user.Active,
		time.Now(),
		time.Now(),
	).Scan(&insertedID)

	if err != nil {
		return 0, fmt.Errorf("error inserting user %w", err)
	}

	return insertedID, nil

}

// func Update() -> updates  user's details based on ID
func (u *UserModel) Update(user *User) error {

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	// this time we are writing to db so we need a statement to execute
	stmt := `update users set 
	email= $1,
	first_name = $2,
	last_name = $3,
	user_active = $4,
	updated_at = $5
	where id = $6
	`
	_, err := u.db.ExecContext(ctx, stmt,
		user.Email,
		user.FirstName,
		user.LastName,
		user.Active,
		time.Now(),
		user.ID,
	)
	if err != nil {
		return fmt.Errorf("error updating user details %w", err)
	}
	return nil
}

// func DeleteById() -> deletes a user based on user ID
func (u *UserModel) DeleteById(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `delete from users where id=$1`

	_, err := u.db.ExecContext(ctx, stmt, id)
	if err != nil {
		return fmt.Errorf("error deleting user %w", err)
	}
	return nil
}

// func GetAll -> returns slice of users sorted by last LastName

func (u *UserModel) GetAll() ([]*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, email, first_name, last_name, password, user_active, created_at, updated_at
						from users order by last_name`

	rows, err := u.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error fetching all users %w", err)
	}
	defer rows.Close()
	var users []*User

	for rows.Next() {
		var u User
		err := rows.Scan(
			&u.ID,
			&u.Email,
			&u.FirstName,
			&u.LastName,
			&u.Password,
			&u.Active,
			&u.CreatedAt,
			&u.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scanning user error %w", err)
		}
		users = append(users, &u)
	}

	return users, nil
}

// func GetByEmail() -> fetches data based on email address

func (u *UserModel) GetByEmail(email string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select id, email, first_name, last_name, password, user_active, created_at, updated_at
						from users where email =$1`

	rows := u.db.QueryRowContext(ctx, query, email)
	var user User

	err := rows.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.Active,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("getting user by email: %w", err)
	}

	return &user, nil

}

// GetByID fetches user by ID
func (u *UserModel) GetByID(id int) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `
		SELECT id, email, first_name, last_name, password, user_active, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var user User
	err := u.db.QueryRowContext(ctx, stmt, id).Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.Active,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("getting user by ID: %w", err)
	}

	return &user, nil
}

// ResetPassword updates the user's password
func (u *UserModel) ResetPassword(id int, password string) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return fmt.Errorf("hashing new password: %w", err)
	}

	stmt := `UPDATE users SET password = $1, updated_at = $2 WHERE id = $3`

	_, err = u.db.ExecContext(ctx, stmt, hashedPassword, time.Now(), id)
	if err != nil {
		return fmt.Errorf("resetting password: %w", err)
	}

	return nil
}

// PasswordMatches compares the plaintext password with the hashed password
func (u *User) PasswordMatches(plainText string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plainText))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return false, nil
		}
		return false, fmt.Errorf("comparing password: %w", err)
	}
	return true, nil
}
