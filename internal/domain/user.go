package domain

import "github.com/google/uuid"

type UserData struct {
	Name string `db:"name" json:"name"`
}
type User struct {
	UserData
	ID uuid.UUID `db:"id" json:"id"`
}
