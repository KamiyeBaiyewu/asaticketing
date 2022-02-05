package errors

import "net/http"

var (

	
	// ErrCreatingContact - Ticket with the ID already exists
	 ErrCreatingContact = APIError{Code: http.StatusInternalServerError, Err: "Error creating contact"}
	

	// ErrPriorityExists - Priority already exists in the database
	// ErrPriorityExists = APIError{Code: http.StatusConflict, Err: "Priority already exists"}

	// ErrCategoryExists - Category already exists in the database
	// ErrCategoryExists = APIError{Code: http.StatusConflict, Err: "Category already exists"}

	// ErrSLAExists - SLA already exists in the database
	// ErrSLAExists = APIError{Code: http.StatusConflict, Err: "SLA already exists"}

	// ErrSourceExists - Source already exists in the database
	// ErrSourceExists = APIError{Code: http.StatusConflict, Err: "Source already exists"}

	// ErrStatusExists - Status already exists in the database
	// ErrStatusExists = APIError{Code: http.StatusConflict, Err: "Status already exists"}

	// ErrTicketExists - User already exists in the database
	// ErrTicketExists = APIError{Code: http.StatusConflict, Err: "Ticket already exists"}

	// ErrTicketIDExists - Ticket with the ID already exists
	// ErrTicketIDExists = APIError{Code: http.StatusConflict, Err: "Ticket with ID already exists"}

	// ErrUpdatingContact - Contact with the ID already exists
	 ErrUpdatingContact = APIError{Code: http.StatusInternalServerError, Err: "Error updating contact"}


	// ErrCauseExists - Cause already exists in the database
	// ErrCauseExists = APIError{Code: http.StatusConflict, Err: "Cause already exists"}
	
)
