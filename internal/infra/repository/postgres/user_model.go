package postgres

import (
	"time"

	"github.com/waliqueiroz/mystery-gifter-api/internal/domain"
)

type UserModel struct {
	ID        string    `db:"id"`
	Name      string    `db:"name"`
	Surname   string    `db:"surname"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func mapUserModelToUser(model UserModel) (*domain.User, error) {
	user := domain.User{
		ID:        model.ID,
		Name:      model.Name,
		Surname:   model.Surname,
		Email:     model.Email,
		Password:  model.Password,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}

	if err := user.Validate(); err != nil {
		return nil, err
	}

	return &user, nil
}
