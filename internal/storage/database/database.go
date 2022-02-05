package database

import (
	"io"

	"github.com/jmoiron/sqlx"
)

// UniqueViolation - Postgres error string for a unique Index violation
const UniqueViolation = "unique_violation"

// Database is a closer interface
type Database interface {
	io.Closer
	ContactsDB
	ClosingRemarkDB
	InboundEmaiiDB //Returns only one email client to connect to
	NoteDB
	ObjectDB
	PolicyDB
	RoleDB
	SessionDB
	UserDB
	UserRoleDB
	SLADB //Service Level Agreement
	// Tickets
	TicketsDB
	TicketCauseDB
	TicketCategoryDB
	TicketPriorityDB
	TicketSourceDB
	TicketStatusDB
	
}

type database struct {
	conn *sqlx.DB
}

func (d *database) Close() error {
	return d.conn.Close()
}
