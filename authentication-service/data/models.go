package data

import "database/sql"

// Models wraps all  models for unified access

type Models struct {
	Users UserModel
}

// NewModels initializes and returns a Models struct

func NewModels(db *sql.DB) Models {
	return Models{
		Users: *NewUserModel(db),
	}
}
