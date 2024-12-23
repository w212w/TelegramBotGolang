package storage

import (
	"TelegramBotGolang/internal/model"
	"context"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/samber/lo"
)

type ArticlePostgresStorage struct {
	db *sqlx.DB
}

type dbArticle struct {
	ID             int64          `db:"a_id"`
	SourcePriority int64          `db:"s_priority"`
	SourceID       int64          `db:"s_id"`
	Title          string         `db:"a_title"`
	Link           string         `db:"a_link"`
	Summary        sql.NullString `db:"a_summary"`
	PublishedAt    time.Time      `db:"a_published_at"`
	PostedAt       sql.NullTime   `db:"a_posted_at"`
	CreatedAt      time.Time      `db:"a_created_at"`
}

func NewArticleStorage(db *sqlx.DB) *ArticlePostgresStorage {
	return &ArticlePostgresStorage{db: db}
}

func (s *ArticlePostgresStorage) Store(ctx context.Context, article model.Article) error {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	if _, err := conn.ExecContext(
		ctx,
		`INSERT INTO articles (source_id, title, link, summary, published at)
			VALUES ($1, $2, $3, $4, $5)
			ON CONFLICT DO NOTHING`,
		article.SourceID,
		article.Title,
		article.Link,
		article.Summary,
		article.PublishedAt,
	); err != nil {
		return err
	}
	return nil

}

func (s *ArticlePostgresStorage) AllNotPosted(ctx context.Context, since time.Time, limit int64) ([]model.Article, error) {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	var articles []dbArticle
	if err := conn.SelectContext(
		ctx,
		&articles,
		`SELECT * FROM articles 
		WHERE posted_at IS NULL 
		AND published_at >= $1::timestamp
		 ORDER BY published_at DESC
		  LIMIT $2`,
		since.UTC().Format(time.RFC3339),
		limit,
	); err != nil {
		return nil, err
	}

	return lo.Map(articles, func(article dbArticle, _ int) model.Article {
		return model.Article{
			ID:          article.ID,
			SourceID:    article.SourceID,
			Title:       article.Title,
			Link:        article.Link,
			Summary:     article.Summary.String,
			PostedAt:    article.PostedAt.Time,
			PublishedAt: article.PublishedAt,
			CreatedAt:   article.CreatedAt,
		}
	}), nil
}

func (s *ArticlePostgresStorage) MarkPosted(ctx context.Context, id int64) error {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	if _, err := conn.ExecContext(
		ctx,
		`UPDATE articles SET posted_at = $1::timestamp WHERE id = $2;`,
		time.Now().UTC().Format(time.RFC3339),
		id,
	); err != nil {
		return err
	}
	return nil
}
