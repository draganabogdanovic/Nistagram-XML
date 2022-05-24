package main

import (
	"context"
	"crypto/rsa"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/KristijanPill/Nishtagram/post-service/handler"
	"github.com/KristijanPill/Nishtagram/post-service/middleware"
	"github.com/KristijanPill/Nishtagram/post-service/model"
	"github.com/KristijanPill/Nishtagram/post-service/repository"
	"github.com/KristijanPill/Nishtagram/post-service/service"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/lestrrat-go/jwx/jwk"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var publicKey *rsa.PublicKey

func init() {
	fetchPublicKey()
}

func initDatabase() *gorm.DB {
	host := os.Getenv("DBHOST")
	user := os.Getenv("USER")
	password := os.Getenv("PASSWORD")
	dbname := os.Getenv("DBNAME")
	dbport := os.Getenv("DBPORT")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai", host, user, password, dbname, dbport)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err.Error())
	}

	db.AutoMigrate(&model.Post{})
	db.AutoMigrate(&model.Comment{})
	db.AutoMigrate(&model.Review{})
	db.AutoMigrate(&model.SavedPost{})
	db.AutoMigrate(&model.Location{})
	db.AutoMigrate(&model.Report{})

	return db
}

func handleFunc(postHandler *handler.PostHandler, commentHandler *handler.CommentHandler, reviewHandler *handler.ReviewHandler,
	savedPostHandler *handler.SavedPostHandler, locationHandler *handler.LocationHandler, securityMiddleware *middleware.SecurityMiddleware, sm *mux.Router) {
	postRouterRestricted := sm.Methods(http.MethodPost).Subrouter()
	postRouterRestricted.HandleFunc("/", postHandler.Create)
	postRouterRestricted.HandleFunc("/save", savedPostHandler.SavePost)
	postRouterRestricted.HandleFunc("/comment", commentHandler.Create)
	postRouterRestricted.HandleFunc("/review", reviewHandler.ReviewPost)
	postRouterRestricted.HandleFunc("/report", postHandler.CreateReport)
	postRouterRestricted.Use(securityMiddleware.Authenticate)

	getRouterPublic := sm.Methods(http.MethodGet).Subrouter()
	getRouterPublic.HandleFunc("/user/{id}", postHandler.FindByOtherUser)
	getRouterPublic.HandleFunc("/location/{query}", locationHandler.GetByQuery)
	getRouterPublic.HandleFunc("/comment/{id}", commentHandler.FindAllByPostID)
	getRouterPublic.HandleFunc("/search/location", postHandler.SearchPostsByLocation)
	getRouterPublic.HandleFunc("/search/tags", postHandler.SearchPostsByTags)
	getRouterPublic.Use(securityMiddleware.UserContext)

	getRouterRestricted := sm.Methods(http.MethodGet).Subrouter()
	getRouterRestricted.HandleFunc("/", postHandler.FindByUser)
	getRouterRestricted.HandleFunc("/liked", postHandler.GetLikedPosts)
	getRouterRestricted.HandleFunc("/disliked", postHandler.GetDislikedPosts)
	getRouterRestricted.HandleFunc("/saved/names", savedPostHandler.GetAllCollectionNames)
	getRouterRestricted.HandleFunc("/saved", savedPostHandler.GetAllByLoggedInUser)
	getRouterRestricted.Use(securityMiddleware.Authenticate)
}

func main() {
	database := initDatabase()

	postRepository := repository.NewPostRepository(database)
	commentRepository := repository.NewCommentRepository(database)
	reviewRepository := repository.NewReviewRepository(database)
	savedPostRepository := repository.NewSavedPostRepository(database)
	locationRepository := repository.NewLocationRepository(database)

	postService := service.NewPostService(postRepository, reviewRepository, commentRepository)
	commentService := service.NewCommentService(commentRepository)
	reviewService := service.NewReviewService(reviewRepository)
	savedPostService := service.NewSavedPostService(savedPostRepository, reviewRepository, commentRepository)
	locationService := service.NewLocationService(locationRepository)

	securityMiddleware := middleware.NewSecurityMiddleware(publicKey)
	postHandler := handler.NewPostHandler(postService)
	commentHandler := handler.NewCommentHandler(commentService)
	reviewHandler := handler.NewReviewHandler(reviewService)
	savedPostHandler := handler.NewSavedPostHandler(savedPostService)
	locationHandler := handler.NewLocationHandler(locationService)

	sm := mux.NewRouter()

	handleFunc(postHandler, commentHandler, reviewHandler, savedPostHandler, locationHandler, securityMiddleware, sm)

	bindAddress := fmt.Sprintf(":%s", os.Getenv("POST_SERVICE_PORT"))

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "content-type", "Authorization", "Referer"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "POST", "HEAD", "PUT", "DELETE", "OPTIONS"})

	cors := handlers.CORS(headersOk, originsOk, methodsOk)

	s := http.Server{
		Addr:         bindAddress,       // configure the bind address
		Handler:      cors(sm),          // set the default handler
		ReadTimeout:  5 * time.Second,   // max time to read request from the client
		WriteTimeout: 10 * time.Second,  // max time to write response to the client
		IdleTimeout:  120 * time.Second, // max time for connections using TCP Keep-Alive
	}

	// fetch public key every 5 minutes
	go func() {
		for {
			fetchPublicKey()
			time.Sleep(5 * time.Minute)
		}
	}()

	// start the server
	go func() {

		err := s.ListenAndServe()
		if err != nil {
			os.Exit(1)
		}
	}()

	// trap sigterm or interupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	// Block until a signal is received.
	sig := <-c
	log.Println("Got signal:", sig)

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(ctx)
}

func fetchPublicKey() {
	requestURL := fmt.Sprintf("http://%s:%s/public-keys", os.Getenv("AUTH_SERVICE_DOMAIN"), os.Getenv("AUTH_SERVICE_PORT"))
	response, err := http.Get(requestURL)

	if err != nil {
		return
	}

	keys, err := io.ReadAll(response.Body)

	if err != nil {
		return
	}

	jwks, err := jwk.Parse(keys)

	if err != nil {
		return
	}

	key, _ := jwks.Get(0)
	var rawKey interface{}

	if err := key.Raw(&rawKey); err != nil {
		return
	}

	publicKey, _ = rawKey.(*rsa.PublicKey)
}
