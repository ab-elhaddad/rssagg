package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ab-elhaddad/rssagg/internal/database"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT is not set")
	}

	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		log.Fatal("DB_URL is not set")
	}
	conn, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	db := database.New(conn)
	go startScraping(db, 10, time.Minute*10)
	apiCfg := &apiConfig{
		DB: db,
	}

	router := chi.NewRouter()

	router.Use(cors.AllowAll().Handler)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	authRouter := chi.NewRouter()
	router.Mount("/auth", authRouter)
	authRouter.Post("/register", apiCfg.handlerRegisterUser)
	authRouter.Post("/login", apiCfg.handlerLoginUser)

	userRouter := chi.NewRouter()
	router.Mount("/users", userRouter)
	userRouter.Get("/me", apiCfg.middlewareAuth(apiCfg.handlerGetUser))

	feedsRouter := chi.NewRouter()
	router.Mount("/feeds", feedsRouter)
	feedsRouter.Post("/", apiCfg.middlewareAuth(apiCfg.handlerFeedsCreate))
	feedsRouter.Get("/", apiCfg.middlewareAuth(apiCfg.handlerAllFeedsGet))
	feedsRouter.Get("/me", apiCfg.middlewareAuth(apiCfg.handlerGetFeedsByUser))
	feedsRouter.Get("/{feedId}", apiCfg.middlewareAuth(apiCfg.handlerGetFeed))

	feedFollowsRouter := chi.NewRouter()
	router.Mount("/feed-follows", feedFollowsRouter)
	feedFollowsRouter.Post("/", apiCfg.middlewareAuth(apiCfg.handlerCreateFeedFollow))
	feedFollowsRouter.Get("/", apiCfg.middlewareAuth(apiCfg.handlerGetFeedFollows))
	feedFollowsRouter.Delete("/{feedFollowID}", apiCfg.middlewareAuth(apiCfg.handlerDeleteFeedFollow))

	postsRouter := chi.NewRouter()
	router.Mount("/posts", postsRouter)
	postsRouter.Get("/", apiCfg.middlewareAuth(apiCfg.handlerGetPostsForUser))

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to RSS Aggregator"))
	})

	router.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("404 Not Found"))
	})

	log.Printf("Server is running on port %v", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
