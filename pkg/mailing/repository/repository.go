package repository

import (
	"context"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"vse.com/4IT428/2023/newsletter/pkg/mailing/models"
)

type MailingRepository interface {
	GetNewsletterWithPost(ctx context.Context, newsletterId string, postId int) (models.NewsletterPost, error)
}

type SQLMailingRepository struct {
	DB *sqlx.DB
}

func (r *SQLMailingRepository) GetNewsletterWithPost(ctx context.Context, newsletterId string, postId int) (models.NewsletterPost, error) {
	var newsletterPost models.NewsletterPost

	query := `SELECT newsletters.user_id, newsletters.name, posts.heading, posts.body, posts.id
           FROM newsletters
           JOIN posts ON newsletters.id = posts.newsletter_id 
           WHERE newsletters.id = $1 AND posts.id = $2;`

	err := r.DB.GetContext(ctx, &newsletterPost, query, newsletterId, postId)

	if err != nil {
		zap.Error(err)
		return models.NewsletterPost{}, err
	}

	return newsletterPost, nil
}
