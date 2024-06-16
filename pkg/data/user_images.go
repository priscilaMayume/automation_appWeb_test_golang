package data

import "time"

// UserImage é o tipo para imagens de perfil de utilizador.
type UserImage struct {
	ID        int       `json:"id"`
	UserID    int    `json:"user_id"`
	FileName  string    `json:"file_name"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}
