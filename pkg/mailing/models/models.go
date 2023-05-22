package models

import (
	"github.com/google/uuid"
)

type NewsletterPost struct {
	AdminID          uuid.UUID `db:"user_id"`
	NameOfNewsletter string    `db:"name"`
	PostID           int       `db:"id"`
	Heading          string    `db:"heading"`
	Body             string    `db:"body"`
}
