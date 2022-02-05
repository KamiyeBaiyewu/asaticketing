package database

import (
	"context"

	"github.com/lib/pq"
	apiErr "github.com/lilkid3/ASA-Ticket/Backend/internal/api/errors"
	"github.com/lilkid3/ASA-Ticket/Backend/internal/model"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var ()

// TicketsDB - holds the interface to the Tickets
type TicketsDB interface {
	CreateTicket(ctx context.Context, ticket *model.Ticket) (err error)
	GetTicketByID(ctx context.Context, ticketID *model.TicketID) (*model.Ticket, error)
	ListAllTickets(ctx context.Context) ([]*model.Ticket, error)
	UpdateTicket(ctx context.Context, ticket *model.Ticket) error
	DeleteTicket(ctx context.Context, ticketID *model.TicketID) (bool, error)

	/* MISC */
	ListAllTicketNotes(ctx context.Context, ticketID *model.TicketID) ([]*model.Note, error)
	DeleteTicketNote(ctx context.Context,ticketID *model.TicketID, noteID *model.NoteID) (bool, error)
	CloseTicket(ctx context.Context, ticketID *model.TicketID) (bool,error)
	ClosingRemark(ctx context.Context, ticketID *model.TicketID) (*model.ClosingRemark, error)
}

func (d *database) GrantTicket(ctx context.Context) error {
	if _, err := d.conn.ExecContext(ctx, "todo"); err != nil {
		if pqError, ok := err.(*pq.Error); ok {
			if pqError.Code.Name() == "unique_violation" && pqError.Constraint == "user_tickets_user" {
				return nil
			}
			logrus.WithFields(logrus.Fields{
				"PQ Code.Name":   pqError.Code.Name(),
				"PQ Constraints": pqError.Constraint,
				"PQ Column":      pqError.Column,
			}).Info()
		}
		return errors.Wrap(err, "could not grant user ticket")
	}
	return nil
}

const createTicketQuery = `
		INSERT INTO tickets (
			 	subject, description, created_by, category_id, status_id, priority_id, source_id, sla_id,  deadline
			)
			VALUES (
				:subject, :description, :created_by,  :category_id,  :status_id,  :priority_id,  :source_id, :sla_id,  :deadline
				)
				RETURNING ticket_id`

func (d *database) CreateTicket(ctx context.Context, ticket *model.Ticket) (err error) {
	rows, err := d.conn.NamedQueryContext(ctx, createTicketQuery, ticket)
	if rows != nil {
		defer rows.Close()
	}

	if err != nil {
		if pqError, ok := err.(*pq.Error); ok {
			if pqError.Constraint == "unique_ticket" {
				err = apiErr.ErrTicketExists
				return
			}
			switch pqError.Code.Name() {
			case "unique_violation":
				if pqError.Constraint == "tickets_pkey" {
					return apiErr.ErrTicketExists
				}
			// One of the Foreign key ID is missing
			case "foreign_key_violation":

				switch pqError.Constraint {
				case "tickets_category_id_fkey":
					return apiErr.ErrNotExist("Category")
				case "tickets_status_id_fkey":
					return apiErr.ErrNotExist("Status")
				case "tickets_priority_id_fkey":
					return apiErr.ErrNotExist("Priority")
				case "tickets_source_id_fkey":
					return apiErr.ErrNotExist("Source")
				case "tickets_created_by_fkey":
					return apiErr.ErrNotExist("User")
				case "tickets_sla_id_fkey":
					return apiErr.ErrNotExist("SLA")
				case "tickets_assigned_to_fkey":
					return apiErr.ErrNotExist("Assignee")

				}
			}

			logrus.WithFields(logrus.Fields{
				"PQ Code.Name":   pqError.Code.Name(),
				"PQ Constraints": pqError.Constraint,
				"PQ Column":      pqError.Column,
				"Error":          err,
			}).Info()
			return apiErr.ErrCreatingTicket
		}
		logrus.WithError(err).Error()
		return apiErr.ErrCreatingTicket
	}

	rows.Next()
	if err := rows.Scan(&ticket.ID); err != nil {
		err = errors.Wrap(err, "Could not get the Ticket ID")
	}
	return
}

const getTicketByIDQuery = `
	SELECT tk.ticket_id, tk.subject, tk.description, tk.created_by, tk.number, 
	tk.category_id, tk.status_id, tk.priority_id, tk.source_id, tk.sla_id,  
	tk.deadline, tk.closed_at, tk.created_at, tk.updated_at, tk.deleted_at
	FROM tickets tk
	WHERE tk.ticket_id = $1
	AND tk.deleted_at IS NULL
`

func (d *database) GetTicketByID(ctx context.Context, ticketID *model.TicketID) (*model.Ticket, error) {
	ticket := model.Ticket{}
	if err := d.conn.GetContext(ctx, &ticket, getTicketByIDQuery, ticketID); err != nil {

		logrus.WithError(err).Error()
		return nil, apiErr.ErrNotFound
	}
	return &ticket, nil
}

const listAllTicketsQuery = `
	SELECT tk.ticket_id, tk.subject, tk.description, tk.created_by,tk.number, 
	tk.category_id, tk.status_id, tk.priority_id, tk.source_id, tk.sla_id,  
	tk.deadline, tk.closed_at, tk.created_at, tk.updated_at, tk.deleted_at
	FROM tickets tk
	WHERE tk.deleted_at IS NULL
`

func (d *database) ListAllTickets(ctx context.Context) ([]*model.Ticket, error) {
	tickets := []*model.Ticket{}
	if err := d.conn.SelectContext(ctx, &tickets, listAllTicketsQuery); err != nil {
		println(err.Error())
		return nil, err
	}
	return tickets, nil
}

const updateTicketQuery = `
		UPDATE tickets
		SET subject = :subject,
		description = :description,
		category_id = :category_id,
		status_id = :status_id,
		priority_id = :priority_id,
		source_id = :source_id,
		updated_at = NOW()
		WHERE ticket_id = :ticket_id
		AND deleted_at is null`

func (d *database) UpdateTicket(ctx context.Context, ticket *model.Ticket) error {
	result, err := d.conn.NamedExecContext(ctx, updateTicketQuery, ticket)
	if err != nil {
		println("PQERROR => ", err.Error())
		if pqError, ok := err.(*pq.Error); ok {
			switch pqError.Code.Name() {

			// One of the Foreign key ID is missing
			case "foreign_key_violation":
				switch pqError.Constraint {
				case "tickets_category_id_fkey":
					return apiErr.ErrNotExist("Category")
				case "tickets_status_id_fkey":
					return apiErr.ErrNotExist("Status")
				case "tickets_priority_id_fkey":
					return apiErr.ErrNotExist("Priority")
				case "tickets_source_id_fkey":
					return apiErr.ErrNotExist("Source")
				case "tickets_created_by_fkey":
					return apiErr.ErrNotExist("User")
				case "tickets_sla_id_fkey":
					return apiErr.ErrNotExist("SLA")
				case "tickets_assigned_to_fkey":
					return apiErr.ErrNotExist("Assignee")

				}
			}

			logrus.WithFields(logrus.Fields{
				"PQ Code.Name":   pqError.Code.Name(),
				"PQ Constraints": pqError.Constraint,
				"PQ Column":      pqError.Column,
			}).Info()
			return apiErr.ErrUpdatingTicket
		}
		return apiErr.ErrUpdatingTicket
	}

	rows, err := result.RowsAffected()

	if err != nil || rows == 0 {
		return errors.New("Ticket Not found")
	}

	return nil
}

const deleteTicketQuery = `
	UPDATE tickets
	SET deleted_at = NOW()
	WHERE ticket_id = $1 AND deleted_at is NULL;
	`

func (d *database) DeleteTicket(ctx context.Context, ticketID *model.TicketID) (bool, error) {
	result, err := d.conn.ExecContext(ctx, deleteTicketQuery, ticketID)
	if err != nil {
		return false, err
	}

	rows, err := result.RowsAffected()
	if err != nil || rows == 0 {
		return false, err
	}
	return true, nil
}


const closeTicketquery = `
	UPDATE tickets
	SET closed_at = NOW()
	WHERE ticket_id = $1 AND deleted_at is NULL;
	`

func (d *database) CloseTicket(ctx context.Context, ticketID *model.TicketID) (bool, error) {
	result, err := d.conn.ExecContext(ctx, closeTicketquery, ticketID)
	if err != nil {
		return false, err
	}

	rows, err := result.RowsAffected()
	if err != nil || rows == 0 {
		return false, err
	}
	return true, nil
}


const listAllTicketNotesQuery = `
	SELECT note_id, note, ticket_id, created_by, created_at, updated_at, deleted_at
	from ticket_notes
	WHERE ticket_id = $1
	AND deleted_at IS NULL
`

func (d *database) ListAllTicketNotes(ctx context.Context, ticketID *model.TicketID) ([]*model.Note, error) {
	notes := []*model.Note{}
	if err := d.conn.SelectContext(ctx, &notes, listAllTicketNotesQuery, ticketID); err != nil {
		logrus.Error(err)
		return nil, err
	}
	return notes, nil
}




const deleteTicketNoteQuery = `
	update ticket_notes
	SET deleted_at = NOW()
	WHERE ticket_id = $1 
	AND note_id = $2 
	AND deleted_at is NULL;
	`

func (d *database) DeleteTicketNote(ctx context.Context,ticketID *model.TicketID, noteID *model.NoteID) (bool, error) {

	result, err := d.conn.ExecContext(ctx, deleteTicketNoteQuery,ticketID, noteID)
	if err != nil {
		return false, err
	}

	rows, err := result.RowsAffected()
	if err != nil || rows == 0 {
		return false, err
	}
	return true, nil
}

const closingRemarkForTicket = `
	SELECT remark_id, ticket_id, cause_id, closed_by, remark, created_at, updated_at, deleted_at
	FROM ticket_closing_remarks
	WHERE ticket_id = $1 
	AND deleted_at is NULL`

func (d *database) ClosingRemark(ctx context.Context, ticketID *model.TicketID) (*model.ClosingRemark, error){
	closingRemark := model.ClosingRemark{}
	if err := d.conn.GetContext(ctx, &closingRemark, closingRemarkForTicket, ticketID); err != nil {
		return nil, apiErr.ErrNotFound
	}
	return &closingRemark, nil


}