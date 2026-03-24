package storage

import (
	"github.com/nicholaskim7/rec-programming/internal/models"
	"github.com/nicholaskim7/rec-programming/internal/utils"
	"database/sql"
	"encoding/hex"
	"context"
	"fmt"
)

// swap out user slice with sqlite db
type UserStore struct {
	db *sql.DB
	hasher *utils.Argon2idHash
}

func NewUserStore(db *sql.DB, hasher *utils.Argon2idHash) *UserStore {
	return &UserStore{
		db: db,
		hasher: hasher,
	}
}

func (s *UserStore) Create(ctx context.Context, user models.User) (models.User, error){
	hashSalt, err := s.hasher.GenerateHash([]byte(user.Password), nil)
	if err != nil {
		return models.User{}, err
	}

	// Convert raw bytes to Hex strings for safe TEXT storage
    hashHex := hex.EncodeToString(hashSalt.Hash)
    saltHex := hex.EncodeToString(hashSalt.Salt)

	query := `INSERT INTO users (first_name, last_name, email, password, salt, username) 
			  VALUES (?, ?, ?, ?, ?, ?)
			  RETURNING id, first_name, last_name, email, username, date_created`

	err = s.db.QueryRowContext(
		ctx,
		query, 
		user.FirstName, 
		user.LastName, 
		user.Email, 
		hashHex,
		saltHex,
		user.Username,
	).Scan(
		&user.ID, 
        &user.FirstName, 
        &user.LastName, 
        &user.Email, 
        &user.Username,
        &user.DateCreated,
	)
	if err != nil {
        return models.User{}, err
    }
	// clear user password before response back to client
	user.Password = ""
	return user, nil
}


func (s *UserStore) Get(ctx context.Context) ([]models.UserResponse, error) {
	// respond without user passwords
	users := []models.UserResponse{}

	query := `SELECT id, first_name, last_name, email, username, date_created FROM users`
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var u models.UserResponse
		// Scan each row and copy the column values into the struct fields
		if err := rows.Scan(
			&u.ID, 
			&u.FirstName,
			&u.LastName,
			&u.Email,
			&u.Username,
			&u.DateCreated,
		); err != nil {
			return users, fmt.Errorf("scan failed: %w", err)
		}
		users = append(users, u)
	}

	if err = rows.Err(); err != nil {
		return users, fmt.Errorf("rows iteration error: %w", err)
	}

	return users, nil
}

// get user by username (include password in this response so we can compare password to hash)
func (s *UserStore) GetByUsername(ctx context.Context, username string) (models.User, error) {
	user := models.User{}

	query := `SELECT id, first_name, last_name, email, username, password, salt, date_created FROM users
			  WHERE username = ?`
	err := s.db.QueryRowContext(
		ctx, 
		query, 
		username,
	).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Username,
		&user.Password,
		&user.Salt,
		&user.DateCreated,
	)
	if err != nil {
        if err == sql.ErrNoRows {
            return models.User{}, fmt.Errorf("user not found")
        }
        return models.User{}, fmt.Errorf("query failed: %w", err)
    }
	return user, nil
}

func (s *UserStore) GetByID(ctx context.Context, id int64) (models.UserResponse, error) {
	user := models.UserResponse{}

	query := `SELECT id, first_name, last_name, email, username, date_created FROM users
			  WHERE id = ?`
	err := s.db.QueryRowContext(
		ctx, 
		query, 
		id,
	).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Username,
		&user.DateCreated,
	)
	if err != nil {
        if err == sql.ErrNoRows {
            return models.UserResponse{}, fmt.Errorf("user not found")
        }
        return models.UserResponse{}, fmt.Errorf("query failed: %w", err)
    }
    return user, nil
}

func (s *UserStore) Login(ctx context.Context, loginPayload models.UserLoginPayload) (models.User, error) {
	// call db method to fetch user by username
	user, err := s.GetByUsername(ctx, loginPayload.Username)
	if err != nil {
		return models.User{}, err
	}

	// decode hex hash and salt that were safely stored into TEXT
	decodedHash, err := hex.DecodeString(user.Password)
	if err != nil {
        return models.User{}, fmt.Errorf("failed to decode hash: %w", err)
    }
	decodedSalt, err := hex.DecodeString(user.Salt)
    if err != nil {
        return models.User{}, fmt.Errorf("failed to decode salt: %w", err)
    }

	err = s.hasher.Compare(decodedHash, decodedSalt, []byte(loginPayload.Password))
	if err != nil {
		// passwords dont match
		return models.User{}, err
	}
	// clear user password, dont send hash to client
	user.Password = ""
	user.Salt = ""

	return user, nil
}

