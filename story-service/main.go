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

	"github.com/KristijanPill/Nishtagram/story-service/handler"
	"github.com/KristijanPill/Nishtagram/story-service/middleware"
	"github.com/KristijanPill/Nishtagram/story-service/model"
	"github.com/KristijanPill/Nishtagram/story-service/repository"
	"github.com/KristijanPill/Nishtagram/story-service/service"
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

	db.AutoMigrate(&model.Story{})
	db.AutoMigrate(&model.StoryHighlight{})
	db.AutoMigrate(&model.Report{})

	return db
}

func handleFunc(storyHandler *handler.StoryHandler, storyHighlightHandler *handler.StoryHighlightHandler, securityMiddleware *middleware.SecurityMiddleware, sm *mux.Router) {
	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/", storyHandler.Create)
	postRouter.HandleFunc("/highlight", storyHighlightHandler.HighlightStory)
	postRouter.HandleFunc("/report", storyHandler.CreateReport)

	postRouter.Use(securityMiddleware.Authenticate)

	getRouterPublic := sm.Methods(http.MethodGet).Subrouter()
	getRouterPublic.HandleFunc("/user/{id}", storyHandler.FindByUser)
	getRouterPublic.Use(securityMiddleware.UserContext)

	getRouterRestricted := sm.Methods(http.MethodGet).Subrouter()
	getRouterRestricted.HandleFunc("/", storyHandler.FindByLoggedInUser)
	getRouterRestricted.HandleFunc("/all", storyHandler.FindAllByLoggedInUser)
	getRouterRestricted.HandleFunc("/highlight/names", storyHighlightHandler.GetAllHighlightNames)
	getRouterRestricted.HandleFunc("/highlight", storyHighlightHandler.GetAllByLoggedInUser)
	getRouterRestricted.Use(securityMiddleware.Authenticate)

}

func main() {
	database := initDatabase()

	storyRepository := repository.NewStoryRepository(database)
	storyHighlightRepository := repository.NewStoryHighlightRepository(database)

	storyService := service.NewStoryService(storyRepository)
	storyHighlightService := service.NewStoryHighlightService(storyHighlightRepository, storyRepository)

	securityMiddleware := middleware.NewSecurityMiddleware(publicKey)
	storyHandler := handler.NewStoryHandler(storyService)
	storyHighlightHandler := handler.NewStoryHighlightHandler(storyHighlightService)

	sm := mux.NewRouter()

	handleFunc(storyHandler, storyHighlightHandler, securityMiddleware, sm)

	bindAddress := fmt.Sprintf(":%s", os.Getenv("STORY_SERVICE_PORT"))

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "content-type", "Authorization"})
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
