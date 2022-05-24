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

	"github.com/KristijanPill/Nishtagram/user-service/handler"
	"github.com/KristijanPill/Nishtagram/user-service/middleware"
	"github.com/KristijanPill/Nishtagram/user-service/model"
	"github.com/KristijanPill/Nishtagram/user-service/repository"
	"github.com/KristijanPill/Nishtagram/user-service/service"
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

	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.Follow{})
	db.AutoMigrate(&model.FollowRequest{})
	db.AutoMigrate(&model.VerificationRequest{})
	db.AutoMigrate(&model.Block{})

	return db
}

func handleFunc(handler *handler.UserHandler, followHandler *handler.FollowHandler, followRequestHandler *handler.FollowRequestHandler,
	verificationRequestHandler *handler.VerificationRequestHandler, blockHandler *handler.BlockHandler, securityMiddleware *middleware.SecurityMiddleware, sm *mux.Router) {
	getRouterRestricted := sm.Methods(http.MethodGet).Subrouter()
	getRouterRestricted.HandleFunc("/user-info", handler.GetUserInfo)
	getRouterRestricted.HandleFunc("/user-profile-info", handler.GetUserProfileInfo)
	getRouterRestricted.Use(securityMiddleware.Authenticate)

	getRouterPublic := sm.Methods(http.MethodGet).Subrouter()
	getRouterPublic.HandleFunc("/follow/outgoing/{id}", followHandler.GetOutgoingFollowStatus)
	getRouterPublic.HandleFunc("/search", handler.SearchUsers)
	getRouterPublic.HandleFunc("/user-profile-info/{username}", handler.GetOtherUserProfileInfo)
	getRouterPublic.Use(securityMiddleware.UserContext)

	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/register", handler.Create)

	postRouterPublic := sm.Methods(http.MethodPost).Subrouter()
	postRouterPublic.HandleFunc("/users-details", handler.GetUsersDetails)
	postRouterPublic.Use(securityMiddleware.UserContext)

	putRouterRestricted := sm.Methods(http.MethodPut).Subrouter()
	putRouterRestricted.HandleFunc("/", handler.Update)
	putRouterRestricted.HandleFunc("/follow/add-close-friend/{id}", followHandler.AddCloseFriend)
	putRouterRestricted.HandleFunc("/follow/remove-close-friend/{id}", followHandler.RemoveCloseFriend)
	putRouterRestricted.HandleFunc("/follow/mute/{id}", followHandler.Mute)
	putRouterRestricted.HandleFunc("/follow/unmute/{id}", followHandler.Unmute)
	putRouterRestricted.HandleFunc("/profile-picture", handler.UpdateProfilePicture)
	putRouterRestricted.Use(securityMiddleware.Authenticate)

	postRouterRestricted := sm.Methods(http.MethodPost).Subrouter()
	postRouterRestricted.HandleFunc("/follow/{id}", followHandler.Follow)
	postRouterRestricted.HandleFunc("/follow/accept/{id}", followRequestHandler.Accept)
	postRouterRestricted.HandleFunc("/follow/decline/{id}", followRequestHandler.Decline)
	postRouterRestricted.HandleFunc("/verify", verificationRequestHandler.Create)
	postRouterRestricted.HandleFunc("/block/{id}", blockHandler.Block)
	postRouterRestricted.HandleFunc("/unblock/{id}", blockHandler.Unblock)
	postRouterRestricted.Use(securityMiddleware.Authenticate)
}

func main() {
	database := initDatabase()

	userRepository := repository.NewUserRepository(database)
	followRepository := repository.NewFollowRepository(database)
	followRequestRepository := repository.NewFollowRequestRepository(database)
	verificationRequestRepository := repository.NewVerificationRequestRepository(database)
	blockRepository := repository.NewBlockRepository(database)

	userService := service.NewUserService(userRepository, followRepository)
	followService := service.NewFollowService(followRepository, followRequestRepository, userRepository)
	followRequestService := service.NewFollowRequestService(followRequestRepository, followRepository)
	verificationRequestService := service.NewVerificationRequestService(verificationRequestRepository)
	blockService := service.NewBlockService(blockRepository, followRepository, followRequestRepository)

	securityMiddleware := middleware.NewSecurityMiddleware(publicKey)
	userHandler := handler.NewUserHandler(userService)
	followHandler := handler.NewFollowHandler(followService)
	followRequestHandler := handler.NewFollowRequestHandler(followRequestService)
	verificationRequestHandler := handler.NewVerificationRequestHandler(verificationRequestService)
	blockHandler := handler.NewBlockHandler(blockService)

	sm := mux.NewRouter()

	handleFunc(userHandler, followHandler, followRequestHandler, verificationRequestHandler, blockHandler, securityMiddleware, sm)

	bindAddress := fmt.Sprintf(":%s", os.Getenv("USER_SERVICE_PORT"))

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
