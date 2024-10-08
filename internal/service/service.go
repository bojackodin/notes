package service

import (
	"context"
	"time"

	"github.com/bojackodin/notes/internal/entity"
	"github.com/bojackodin/notes/internal/repository"
	"github.com/bojackodin/notes/internal/yandex/speller"
)

type Auth interface {
	CreateUser(ctx context.Context, username, password string) (int64, error)
	GenerateToken(ctx context.Context, username, password string) (string, error)
	ParseToken(token string) (int64, error)
}

type Note interface {
	CreateNote(ctx context.Context, title string, userID int64) (int64, error)
	ListNotes(ctx context.Context, userID int64) ([]entity.Note, error)
}

type Services struct {
	Auth Auth
	Note Note
}

type ServicesDependencies struct {
	Repositories *repository.Repositories
	Speller      speller.Speller

	Secret   string
	TokenTTL time.Duration
}

func NewServices(deps ServicesDependencies) *Services {
	return &Services{
		Auth: NewAuthService(deps.Repositories.User, deps.Secret, deps.TokenTTL),
		Note: NewNoteService(deps.Repositories.Note, deps.Speller),
	}
}
