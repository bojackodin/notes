package postgress

import (
	"context"
	"database/sql"
	"errors"

	"github.com/bojackodin/notes/internal/entity"
	"github.com/bojackodin/notes/internal/repository/repositoryerror"
	"github.com/lib/pq"
)

type UserRepository struct {
	client *sql.DB
}

func NewUserRepository(client *sql.DB) *UserRepository {
	return &UserRepository{
		client: client,
	}
}

func (db *UserRepository) CreateUser(ctx context.Context, user entity.User) (int64, error) {
	query := `
		INSERT INTO users (username, password) 
        VALUES ($1, $2)
        RETURNING id`

	err := db.client.QueryRowContext(ctx, query, user.Username, user.Password).Scan(&user.ID)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return 0, repositoryerror.ErrDuplicate
		}
		return 0, err
	}

	return user.ID, nil
}

func (db *UserRepository) GetUserByUsername(ctx context.Context, username string) (entity.User, error) {
	query := `
		SELECT id, username, password, created_at
        FROM users
        WHERE username = $1
        `

	var user entity.User

	err := db.client.QueryRowContext(ctx, query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.CreatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return entity.User{}, repositoryerror.ErrRecordNotFound
		default:
			return entity.User{}, err
		}
	}

	return user, nil
}

func (db *UserRepository) GetUserById(ctx context.Context, id int64) (entity.User, error) {
	return entity.User{}, nil
}
