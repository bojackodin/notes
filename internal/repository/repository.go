package repository

import (
	"context"
	"database/sql"

	"github.com/bojackodin/notes/internal/entity"
	"github.com/bojackodin/notes/internal/repository/postgress"
)

type User interface {
	CreateUser(ctx context.Context, user entity.User) (int64, error)
	GetUserByUsername(ctx context.Context, username string) (entity.User, error)
	GetUserById(ctx context.Context, id int64) (entity.User, error)
}

type Note interface {
	CreateNote(ctx context.Context, note *entity.Note) error
	ListNotes(ctx context.Context, userID int64) ([]*entity.Note, error)
}

type Repositories struct {
	User
	Note
}

func NewRepositories(client *sql.DB) *Repositories {
	return &Repositories{
		User: postgress.NewUserRepository(client),
		Note: postgress.NewNoteRepository(client),
	}
}
