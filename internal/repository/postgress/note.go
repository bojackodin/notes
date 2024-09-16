package postgress

import (
	"context"
	"database/sql"

	"github.com/bojackodin/notes/internal/entity"
)

type NoteRepository struct {
	client *sql.DB
}

func NewNoteRepository(client *sql.DB) *NoteRepository {
	return &NoteRepository{
		client: client,
	}
}

func (db *NoteRepository) CreateNote(ctx context.Context, note *entity.Note) error {
	query := `
		INSERT INTO notes (user_id, title)
		VALUES ($1, $2)
		RETURNING id`

	return db.client.QueryRowContext(ctx, query, note.UserID, note.Title).Scan(&note.ID)
}

func (db *NoteRepository) ListNotes(ctx context.Context, userID int64) ([]*entity.Note, error) {
	query := `
		SELECT id, user_id, title
		FROM notes
		WHERE user_id = $1`

	rows, err := db.client.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	notes := []*entity.Note{}

	for rows.Next() {
		var note entity.Note

		err := rows.Scan(
			&note.ID,
			&note.UserID,
			&note.Title,
		)
		if err != nil {
			return nil, err
		}

		notes = append(notes, &note)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return notes, nil
}
