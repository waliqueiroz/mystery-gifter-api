package domain

//go:generate go run go.uber.org/mock/mockgen -destination mock_domain/identity_generator.go . IdentityGenerator

type IdentityGenerator interface {
	Generate() (string, error)
}
