package service

import (
	"context"
	"errors"

	"github.com/bojackodin/notes/internal/entity"
	"github.com/bojackodin/notes/internal/repository"
	"github.com/bojackodin/notes/internal/service/serviceerror"
	"github.com/bojackodin/notes/internal/yandex/speller"
)

type NoteService struct {
	noteRepository repository.Note
	speller        speller.Speller
}

func NewNoteService(noteRepository repository.Note, speller speller.Speller) *NoteService {
	return &NoteService{
		noteRepository: noteRepository,
		speller:        speller,
	}
}

func (s *NoteService) CreateNote(ctx context.Context, title string, userID int64) (int64, error) {
	err := s.speller.Check(ctx, title)
	if err != nil {
		if errors.Is(err, speller.ErrorSpell{}) {
			return 0, serviceerror.ErrSpeller
		}
		return 0, err
	}

	note := entity.Note{
		Title:  title,
		UserID: userID,
	}

	err = s.noteRepository.CreateNote(ctx, &note)
	if err != nil {
		return 0, err
	}

	return note.ID, nil
}

func (s *NoteService) ListNotes(ctx context.Context, userID int64) ([]*entity.Note, error) {
	return s.noteRepository.ListNotes(ctx, userID)
}
