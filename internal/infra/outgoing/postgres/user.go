package postgres

import (
	"time"

	"github.com/waliqueiroz/mystery-gifter-api/internal/domain"
)

type User struct {
	ID        string    `db:"id"`
	Name      string    `db:"name"`
	Surname   string    `db:"surname"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func mapUserToDomain(user User) (*domain.User, error) {
	domainUser := domain.User{
		ID:        user.ID,
		Name:      user.Name,
		Surname:   user.Surname,
		Email:     user.Email,
		Password:  user.Password,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	if err := domainUser.Validate(); err != nil {
		return nil, err
	}

	return &domainUser, nil
}

func mapUsersToDomain(users []User) ([]domain.User, error) {
	domainUsers := make([]domain.User, 0, len(users))

	for _, model := range users {
		user, err := mapUserToDomain(model)
		if err != nil {
			return nil, err
		}

		domainUsers = append(domainUsers, *user)
	}

	return domainUsers, nil
}
