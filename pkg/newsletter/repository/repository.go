package repository

import (
	"context"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"vse.com/4IT428/2023/newsletter/pkg/newsletter/models"
)

type NewsletterRepository interface {
	GetNewsletters(ctx context.Context) ([]models.Newsletter, error)
	CreateNewsletter(ctx context.Context, newsletter models.Newsletter) (models.Newsletter, error)
	CreateNewsletterPost(ctx context.Context, post models.Post) (models.Post, error)
	GetNewsletterById(ctx context.Context, id string) (models.Newsletter, error)
	UpdateNewsletterById(ctx context.Context, newsletter models.Newsletter) (models.Newsletter, error)
	DeletePostsByNewsletterId(ctx context.Context, newsletterId string) error
	DeleteNewsletterById(ctx context.Context, id string) error
}

type SQLNewsletterRepository struct {
	DB *sqlx.DB
}

func (r *SQLNewsletterRepository) GetNewsletters(ctx context.Context) ([]models.Newsletter, error) {
	query := `SELECT * FROM newsletters`

	var newsletters []models.Newsletter
	err := r.DB.SelectContext(ctx, &newsletters, query)
	if err != nil {
		zap.Error(err)
		return nil, err
	}

	return newsletters, nil
}

func (r *SQLNewsletterRepository) CreateNewsletter(ctx context.Context, newsletter models.Newsletter) (models.Newsletter, error) {
	query := `INSERT INTO newsletters (id, user_id, name)
	          VALUES ($1, $2, $3)`

	_, err := r.DB.ExecContext(ctx, query, newsletter.ID, newsletter.UserId, newsletter.Name)
	if err != nil {
		zap.Error(err)
		return models.Newsletter{}, err
	}

	return newsletter, nil
}

func (r *SQLNewsletterRepository) CreateNewsletterPost(ctx context.Context, post models.Post) (models.Post, error) {
	query := `INSERT INTO posts (newsletter_id, heading, body)
	          VALUES ($1, $2, $3)`

	_, err := r.DB.ExecContext(ctx, query, post.NewsletterID, post.Heading, post.Body)
	if err != nil {
		zap.Error(err)
		return models.Post{}, err
	}

	return post, nil
}

func (r *SQLNewsletterRepository) GetNewsletterById(ctx context.Context, id string) (models.Newsletter, error) {
	query := `SELECT * FROM newsletters WHERE id=$1`

	var newsletter models.Newsletter
	err := r.DB.GetContext(ctx, &newsletter, query, id)
	if err != nil {
		zap.Error(err)
		return models.Newsletter{}, err
	}

	return newsletter, nil
}

func (r *SQLNewsletterRepository) UpdateNewsletterById(ctx context.Context, newsletter models.Newsletter) (models.Newsletter, error) {
	query := `UPDATE newsletters SET name=$1 WHERE id=$2`

	_, err := r.DB.ExecContext(ctx, query, newsletter.Name, newsletter.ID)
	if err != nil {
		zap.Error(err)
		return models.Newsletter{}, err
	}

	return newsletter, nil
}

func (r *SQLNewsletterRepository) DeletePostsByNewsletterId(ctx context.Context, newsletterId string) error {
	query := `DELETE FROM posts WHERE newsletter_id=$1`

	_, err := r.DB.ExecContext(ctx, query, newsletterId)
	if err != nil {
		zap.Error(err)
		return err
	}

	return nil
}

func (r *SQLNewsletterRepository) DeleteNewsletterById(ctx context.Context, newsletterId string) error {
	query := `DELETE FROM newsletters WHERE id=$1`

	_, err := r.DB.ExecContext(ctx, query, newsletterId)
	if err != nil {
		zap.Error(err)
		return err
	}

	return nil
}
