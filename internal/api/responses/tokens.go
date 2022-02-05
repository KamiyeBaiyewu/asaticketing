package responses

import (
	"github.com/lilkid3/ASA-Ticket/Backend/internal/api/auth"
	"github.com/lilkid3/ASA-Ticket/Backend/internal/model"
)

// TokenResponse - structure represents the token information send to the user
type TokenResponse struct {
	Tokens *auth.Tokens `json:"tokens,omitempty"`
	User   *model.User  `json:"user,omitempty"`
}
