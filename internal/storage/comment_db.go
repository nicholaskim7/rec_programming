package storage

import (
	"github.com/nicholaskim7/rec-programming/internal/models"
	"database/sql"
	"context"
	"fmt"
)

type CommentStore struct {
	db *sql.DB
}

func NewCommentStore(db *sql.DB) *CommentStore {
	return &CommentStore{
		db: db,
	}
}

func (s *CommentStore) Create(ctx context.Context, comment models.Comment) (models.Comment, error) {
	query := `INSERT INTO comments (user_id, post_id, comment)
			  VALUES (?, ?, ?)
			  RETURNING id, user_id, post_id, comment, date_created`
	err := s.db.QueryRowContext(
		ctx,
		query, 
		comment.UserID,
		comment.PostID,
		comment.Comment,
	).Scan(
		&comment.ID,
		&comment.UserID,
		&comment.PostID,
		&comment.Comment,
		&comment.DateCreated,
	)
	if err != nil {
		return models.Comment{}, err
	}
	return comment, nil
}

func (s *CommentStore) Get(ctx context.Context) ([]models.Comment, error) {
	comments := []models.Comment{}

	query := `SELECT c.id, c.user_id, u.username, c.post_id, c.comment, c.date_created 
			  FROM comments c
			  INNER JOIN users u ON u.id = c.user_id
	`
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var c models.Comment
		// Scan each row and copy the column values into the struct fields
		if err := rows.Scan(
			&c.ID, 
			&c.UserID,
			&c.Username,
			&c.PostID,
			&c.Comment,
			&c.DateCreated,
		); err != nil {
			return comments, fmt.Errorf("scan failed: %w", err)
		}
		comments = append(comments, c)
	}
	if err = rows.Err(); err != nil {
		return comments, fmt.Errorf("rows iteration error: %w", err)
	}
	return comments, nil
}


func (s *CommentStore) GetByPost(ctx context.Context, postID int64) ([]models.Comment, error) {
	comments := []models.Comment{}

	query := `SELECT c.id, c.user_id, u.username, c.post_id, c.comment, c.date_created
          FROM comments c
          INNER JOIN users u ON u.id = c.user_id
          WHERE c.post_id = ?`
	rows, err := s.db.QueryContext(ctx, query, postID)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var c models.Comment
		// Scan each row and copy the column values into the struct fields
		if err := rows.Scan(
			&c.ID, 
			&c.UserID,
			&c.Username,
			&c.PostID,
			&c.Comment,
			&c.DateCreated,
		); err != nil {
			return comments, fmt.Errorf("scan failed: %w", err)
		}
		comments = append(comments, c)
	}

	if err = rows.Err(); err != nil {
		return comments, fmt.Errorf("rows iteration error: %w", err)
	}
	return comments, nil
}

