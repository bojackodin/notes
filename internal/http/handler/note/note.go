package note

import (
	"errors"
	"net/http"

	"github.com/bojackodin/notes/internal/http/encoding"
	contexthelper "github.com/bojackodin/notes/internal/http/handler/context"
	"github.com/bojackodin/notes/internal/http/httperror"
	"github.com/bojackodin/notes/internal/log"
	"github.com/bojackodin/notes/internal/service"
	"github.com/bojackodin/notes/internal/service/serviceerror"
)

type Controller struct {
	notes service.Note
}

func New(notes service.Note) *Controller {
	return &Controller{
		notes: notes,
	}
}

type createNoteInput struct {
	Title string `json:"title"`
}

type createNoteResponse struct {
	ID int64 `json:"id"`
}

func (ctrl *Controller) CreateNote(w http.ResponseWriter, r *http.Request) error {
	logger := log.FromContext(r.Context())
	userID := contexthelper.ContextGetUserID(r)

	var input createNoteInput
	if err := encoding.Decode(r, &input); err != nil {
		logger.Error("failed to decode body", log.Err(err))
		return httperror.WithStatusError(err, http.StatusBadRequest)
	}

	id, err := ctrl.notes.CreateNote(r.Context(), input.Title, userID)
	if err != nil {
		code := http.StatusInternalServerError
		if errors.Is(err, serviceerror.ErrSpeller) {
			code = http.StatusUnprocessableEntity
		}
		logger.Error("failed to create task", log.Err(err))
		return httperror.WithStatusError(err, code)
	}

	_ = encoding.Encode(http.StatusCreated, w, &createNoteResponse{ID: id})
	return nil
}

type noteResponse struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
}

type listNotesResponse []*noteResponse

func (ctrl *Controller) ListNotes(w http.ResponseWriter, r *http.Request) error {
	logger := log.FromContext(r.Context())
	userID := contexthelper.ContextGetUserID(r)

	notes, err := ctrl.notes.ListNotes(r.Context(), userID)
	if err != nil {
		logger.Error("failed to list tasks", log.Err(err))
		return err
	}

	response := make(listNotesResponse, 0, len(notes))
	for _, note := range notes {
		response = append(response, &noteResponse{
			ID:    note.ID,
			Title: note.Title,
		})
	}

	_ = encoding.Encode(http.StatusOK, w, &response)

	return nil
}
