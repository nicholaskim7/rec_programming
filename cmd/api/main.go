package main
import (
	"net/http"
	"github.com/nicholaskim7/rec-programming/internal/storage"
	"github.com/nicholaskim7/rec-programming/internal/handlers"
	"github.com/nicholaskim7/rec-programming/internal/utils"
	"github.com/nicholaskim7/rec-programming/internal/middleware"
	"fmt"
	"github.com/joho/godotenv"
	"database/sql"
	_ "github.com/glebarez/go-sqlite"
	"github.com/rs/cors"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found")
	}

	// Connect to the SQLite database
    db, err := sql.Open("sqlite", "./my.db?_pragma=foreign_keys(1)")
    if err != nil {
        fmt.Println(err)
        return
    }
	defer db.Close()
	if err := db.Ping(); err != nil {
        fmt.Printf("Error connecting to database: %v\n", err)
        return
    }
    fmt.Println("Connected to the SQLite database successfully.")

	statements := `
	CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY,
        first_name     TEXT NOT NULL,
        last_name TEXT NOT NULL,
		email TEXT NOT NULL,
		password TEXT NOT NULL,
		salt TEXT NOT NULL,
		username TEXT NOT NULL,
		date_created DATETIME DEFAULT CURRENT_TIMESTAMP
    );
	CREATE TABLE IF NOT EXISTS posts (
        id INTEGER PRIMARY KEY,
        user_id INTEGER NOT NULL REFERENCES users(id),
        title TEXT NOT NULL,
		body TEXT NOT NULL,
		tag TEXT,
		date_created DATETIME DEFAULT CURRENT_TIMESTAMP
    );
	CREATE TABLE IF NOT EXISTS comments (
        id INTEGER PRIMARY KEY,
        user_id INTEGER NOT NULL REFERENCES users(id),
		post_id INTEGER NOT NULL REFERENCES posts(id),
		comment TEXT NOT NULL,
		date_created DATETIME DEFAULT CURRENT_TIMESTAMP
    );
	`

    _, err = db.Exec(statements)
	if err != nil {
		fmt.Printf("Failed to create tables: %v\n", err)
        return
	}
	fmt.Println("Created tables")

	// 64 MB memory, 3 iterations, 4 threads, 32-byte hash, 16-byte salt
	hasher := utils.NewArgon2idHash(3, 16, 65536, 4, 32)

	userStore := storage.NewUserStore(db, hasher)
	userHandler := handlers.NewUserHandler(userStore)

	postStore := storage.NewPostStore(db)
	postHandler := handlers.NewPostHandler(postStore)

	commentStore := storage.NewCommentStore(db)
	commentHandler := handlers.NewCommentHandler(commentStore)

	// public routes
	http.HandleFunc("POST /api/users", userHandler.CreateUserHandler)
	http.HandleFunc("GET /api/posts", postHandler.GetPostsHandler)

	http.HandleFunc("GET /api/posts/user/{username}", postHandler.GetPostsByUsernameHandler)
	http.HandleFunc("GET /api/posts/post/{id}", postHandler.GetPostByPostIDHandler)
	http.HandleFunc("GET /api/comments", commentHandler.GetCommentsHandler)
	http.HandleFunc("GET /api/comments/post/{id}", commentHandler.GetCommentsByPostHandler)
	http.HandleFunc("POST /api/users/login", userHandler.LoginHandler)

	// protected routes
	http.HandleFunc("GET /api/users", middleware.AuthMiddleware(userHandler.GetUsersHandler))
	http.HandleFunc("POST /api/posts", middleware.AuthMiddleware(postHandler.CreatePostHandler))
	http.HandleFunc("POST /api/comments", middleware.AuthMiddleware(commentHandler.CreateCommentHandler))
	http.HandleFunc("POST /api/users/logout", middleware.AuthMiddleware(userHandler.LogoutHandler))
	http.HandleFunc("GET /api/users/me", middleware.AuthMiddleware(userHandler.GetMeHandler))
	// serve static html
	// http.HandleFunc("/privacy", func(w http.ResponseWriter, r *http.Request){
	// 	http.ServeFile(w, r, "./static/privacy_policy.html")
	// })

	fmt.Println("starting server on :8080...")
	c := cors.New(cors.Options{
        // Change these to match the URL/port your React app is running on
        AllowedOrigins: []string{"http://localhost:3000", "http://localhost:5173"}, 
        AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowedHeaders: []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "Authorization"},
        AllowCredentials: true, 
    })
	handler := c.Handler(http.DefaultServeMux)

	err = http.ListenAndServe(":8080", handler)
	if err != nil {
		panic(err)
	}
}