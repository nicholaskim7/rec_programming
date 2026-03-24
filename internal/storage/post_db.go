package storage

import (
	"github.com/nicholaskim7/rec-programming/internal/models"
	"database/sql"
	"context"
	"fmt"
)

type PostStore struct {
	db *sql.DB
}

func NewPostStore(db *sql.DB) *PostStore {
	return &PostStore{
		db: db,
	}
}

func (s *PostStore) Create(ctx context.Context, post models.Post) (models.Post, error){
	query := `INSERT INTO posts (user_id, title, body, tag) 
			  VALUES (?, ?, ?, ?)
			  RETURNING id, user_id, title, body, tag, date_created`

	err := s.db.QueryRowContext(
		ctx,
		query, 
		post.UserID, 
		post.Title, 
		post.Body, 
		post.Tag, 
	).Scan(
		&post.ID, 
        &post.UserID, 
        &post.Title, 
        &post.Body, 
        &post.Tag,
        &post.DateCreated,
	)
	if err != nil {
        return models.Post{}, err
    }
	return post, nil
}


func (s *PostStore) Get(ctx context.Context) ([]models.Post, error) {
	posts := []models.Post{}

	query := `
		SELECT p.id, p.user_id, p.title, p.body, p.tag, p.date_created, u.username 
		FROM posts p
		INNER JOIN users u ON u.id = p.user_id
	`
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var p models.Post
		// Scan each row and copy the column values into the struct fields
		if err := rows.Scan(
			&p.ID, 
			&p.UserID,
			&p.Title,
			&p.Body,
			&p.Tag,
			&p.DateCreated,
			&p.Username,
		); err != nil {
			return posts, fmt.Errorf("scan failed: %w", err)
		}
		posts = append(posts, p)
	}

	if err = rows.Err(); err != nil {
		return posts, fmt.Errorf("rows iteration error: %w", err)
	}

	return posts, nil
}



func (s *PostStore) GetByUsername(ctx context.Context, username string) ([]models.Post, error) {
	posts := []models.Post{}

	query := `SELECT p.id, p.user_id, p.title, p.body, p.tag, p.date_created, u.username
			  FROM posts p
			  INNER JOIN users u ON u.id = p.user_id
			  WHERE u.username = ?
			  `
	rows, err := s.db.QueryContext(ctx, query, username)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var p models.Post
		// Scan each row and copy the column values into the struct fields
		if err := rows.Scan(
			&p.ID, 
			&p.UserID,
			&p.Title,
			&p.Body,
			&p.Tag,
			&p.DateCreated,
			&p.Username,
		); err != nil {
			return posts, fmt.Errorf("scan failed: %w", err)
		}
		posts = append(posts, p)
	}

	if err = rows.Err(); err != nil {
		return posts, fmt.Errorf("rows iteration error: %w", err)
	}

	return posts, nil
}


func (s *PostStore) GetByPostID(ctx context.Context, postID int64) (models.Post, error) {
	post := models.Post{}

	query := `SELECT p.id, p.user_id, p.title, p.body, p.tag, p.date_created, u.username 
			  FROM posts p
			  INNER JOIN users u ON u.id = p.user_id 
			  WHERE p.id = ?
			  `
	err := s.db.QueryRowContext(
		ctx, 
		query, 
		postID,
	).Scan(
		&post.ID,
		&post.UserID,
		&post.Title,
		&post.Body,
		&post.Tag,
		&post.DateCreated,
		&post.Username,
	)
	if err != nil {
		return models.Post{}, fmt.Errorf("query failed: %w", err)
	}


	return post, nil
}