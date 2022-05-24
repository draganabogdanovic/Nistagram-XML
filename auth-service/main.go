package main

import (
	"context"
	"crypto/rsa"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/KristijanPill/Nishtagram/auth-service/handler"
	"github.com/KristijanPill/Nishtagram/auth-service/middleware"
	"github.com/KristijanPill/Nishtagram/auth-service/payload"
	"github.com/KristijanPill/Nishtagram/auth-service/repository"
	"github.com/KristijanPill/Nishtagram/auth-service/service"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	privKeyPath = "./keys/app.rsa"
	pubKeyPath  = "./keys/app.rsa.pub"
	hmacKeyPath = "./keys/app.hmac"
)

var (
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
	hmacKey    []byte
)

func init() {
	signBytes, err := ioutil.ReadFile(privKeyPath)
	fatal(err)

	privateKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	fatal(err)

	verifyBytes, err := ioutil.ReadFile(pubKeyPath)
	fatal(err)

	publicKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	fatal(err)

	hmacKey, err = ioutil.ReadFile(hmacKeyPath)
	fatal(err)
}

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
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

	db.AutoMigrate(&payload.Credentials{})
	db.AutoMigrate(&payload.RefreshToken{})

	return db
}

func initRefreshTokenRepository(db *gorm.DB) *repository.RefreshTokenRepository {
	return repository.NewRefreshTokenRepository(db)
}

func initCredentialsRepository(db *gorm.DB) *repository.CredentialsRepository {
	return repository.NewCredentialsRepository(db)
}

func initAuthService(refreshTokenRepository *repository.RefreshTokenRepository, credentialsRepository *repository.CredentialsRepository) *service.AuthService {
	return service.NewAuthService(publicKey, privateKey, hmacKey, refreshTokenRepository, credentialsRepository)
}

func initAuthHandler(service *service.AuthService) *handler.AuthHandler {
	return handler.NewAuthHandler(service)
}

func handleFunc(handler *handler.AuthHandler, refreshMiddleware *middleware.RefreshMiddleware, sm *mux.Router) {
	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/login", handler.Login)
	postRouter.HandleFunc("/register", handler.Register)

	postRouterRestricted := sm.Methods(http.MethodPost).Subrouter()
	postRouterRestricted.Use(refreshMiddleware.Authenticate)
	postRouterRestricted.HandleFunc("/refresh", handler.Refresh)

	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/public-keys", handler.GetPublicKeys)
}

func main() {
	database := initDatabase()
	refreshTokenRepository := initRefreshTokenRepository(database)
	credentialsRepository := initCredentialsRepository(database)
	authService := initAuthService(refreshTokenRepository, credentialsRepository)
	authHandler := initAuthHandler(authService)

	refreshMiddleware := middleware.NewRefreshMiddleware(publicKey, hmacKey)

	sm := mux.NewRouter()

	handleFunc(authHandler, refreshMiddleware, sm)

	bindAddress := fmt.Sprintf(":%s", os.Getenv("AUTH_SERVICE_PORT"))

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "content-type"})
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
