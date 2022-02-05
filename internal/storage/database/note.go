package database

import (
	"context"

	"github.com/lib/pq"
	"github.com/lilkid3/ASA-Ticket/Backend/internal/model"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	apiErr  "github.com/lilkid3/ASA-Ticket/Backend/internal/api/errors"
)

var (
	
)

// NoteDB - Interface holds all the methods for storing and retrieving notes
type NoteDB interface {
	CreateNote(ctx context.Context, note *model.Note) error
	GetNoteByID(ctx context.Context, noteID *model.NoteID) (*model.Note, error)

	UpdateNote(ctx context.Context, note *model.Note) error
	DeleteNote(ctx context.Context, noteID *model.NoteID) (bool, error)
	ListAllNotes(ctx context.Context) ([]*model.Note, error)


}

const createNoteQuery = `
		INSERT INTO ticket_notes (
			note, ticket_id, created_by
			)
			VALUES (
				 :note, :ticket_id, :created_by
				)
				RETURNING note_id`

func (d *database) CreateNote(ctx context.Context, userNote *model.Note) (err error) {

	rows, err := d.conn.NamedQueryContext(ctx, createNoteQuery, userNote)
	if rows != nil {
		defer rows.Close()
	}

	if err != nil {
		if pqError, ok := err.(*pq.Error); ok {
			logrus.WithFields(logrus.Fields{
				"PQ Code.Name":   pqError.Code.Name(),
				"PQ Constraints": pqError.Constraint,
				"PQ Column":      pqError.Column,
			}).Info()
		}
		return
	}

	rows.Next()
	if err := rows.Scan(&userNote.ID); err != nil {
		err = errors.Wrap(err, "Could not get the Note ID")
	}
	return
}

const getNoteByIDQuery = `
	SELECT note_id, note, ticket_id, created_by, created_at, updated_at, deleted_at
	from ticket_notes
	WHERE note_id = $1
	AND deleted_at IS NULL`

func (d *database) GetNoteByID(ctx context.Context, noteID *model.NoteID) (*model.Note, error) {
	userNote := model.Note{}
	if err := d.conn.GetContext(ctx, &userNote, getNoteByIDQuery, noteID); err != nil {

		return nil, apiErr.ErrNoteNotExist
	}
	return &userNote, nil
}

const updateNoteQuery = `
		update ticket_notes
		SET note = :note,
		updated_at = NOW()
		WHERE note_id = :note_id`

func (d *database) UpdateNote(ctx context.Context, note *model.Note) error {
	result, err := d.conn.NamedExecContext(ctx, updateNoteQuery, note)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil || rows == 0 {
		return errors.New("Note Not found")
	}

	return nil
}

const deleteNoteQuery = `
	update ticket_notes
	SET deleted_at = NOW()
	WHERE note_id = $1 AND deleted_at is NULL;
	`

func (d *database) DeleteNote(ctx context.Context, noteID *model.NoteID) (bool, error) {

	result, err := d.conn.ExecContext(ctx, deleteNoteQuery, noteID)
	if err != nil {
		return false, err
	}

	rows, err := result.RowsAffected()
	if err != nil || rows == 0 {
		return false, err
	}
	return true, nil
}

const listAllNotesQuery = `
	SELECT note_id, note, ticket_id, created_by, created_at, updated_at, deleted_at
	from ticket_notes
	WHERE deleted_at is NULL;
`

func (d *database) ListAllNotes(ctx context.Context) ([]*model.Note, error) {
	userNotes := []*model.Note{}
	if err := d.conn.SelectContext(ctx, &userNotes, listAllNotesQuery); err != nil {
		return nil, errors.Wrap(err, "could not get notes")
	}
	return userNotes, nil
}
