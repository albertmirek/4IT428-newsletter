package models

import "github.com/google/uuid"

type Newsletter struct {
	ID     uuid.UUID `db:"id"`
	UserId uuid.UUID `db:"user_id"`
	Name   string    `db:"name"`
}
type CreateNewsletterRequest struct {
	Name string `json:"name"`
}
type CreateNewsletterPost struct {
	Heading string `json:"heading"`
	Body    string `json:"body"`
}
type Post struct {
	ID           uuid.UUID `db:"id"`
	NewsletterID uuid.UUID `db:"newsletterID"`
	Heading      string    `db:"heading"`
	Body         string    `db:"body"`
}

type UpdateNewsletterRequest struct {
	Name string `json:"name"`
}
