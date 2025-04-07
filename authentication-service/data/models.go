package data

import (
	"context"
	"database/sql"
	"time"
)

const dbTimeoutinSec = 3 * time.Second

var db *sql.DB

type Models struct {
	User User
}

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

func New(dbPool *sql.DB) Models {
	db = dbPool
	return Models{
		User: User{},
	}
}

// func GetAll -> returns slice of users sorted by last LastName

func (u *User) GetAll() ([]*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeoutinSec)
	defer cancel()

	query := `select id, email, first_name, last_name, password, user_active, created_at, updated_at
						from users order by last_name`

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []*User

	for rows.Next() {
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
			return nil, err
		}
		users = append(users, &user)
	}

	return users, nil
}

// func GetByEmail() -> fetches data based on email address

func (u *User) GetByEmail(email string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeoutinSec)
	defer cancel()

	query := `select id, email, first_name, last_name, password, user_active, created_at, updated_at
						from users where email =$1`

	rows := db.QueryRowContext(ctx, query, email)
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
		return nil, err
	}

	return &user, nil

}

// func GetById() -> fetches data based on user's id
func (u *User) GetById(id string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeoutinSec)
	defer cancel()

	var user User
	query := `select id, email, first_name, last_name, password, user_active, created_at, updated_at
            from users where id =$1`
	rows := db.QueryRowContext(ctx, query, id)

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
		return nil, err
	}
	return &user, nil

}

// func Update() -> updates one user's details in db
func (u *User) Update() {

}
