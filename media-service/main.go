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

	"github.com/KristijanPill/Nishtagram/media-service/handler"
	"github.com/KristijanPill/Nishtagram/media-service/middleware"
	"github.com/KristijanPill/Nishtagram/media-service/service"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/lestrrat-go/jwx/jwk"
)

var publicKey *rsa.PublicKey

func init() {
	fetchPublicKey()
}

func handleFunc(handler *handler.MediaHandler, securityMiddleware *middleware.SecurityMiddleware, sm *mux.Router, fs *http.Handler) {
	postRouterRestricted := sm.Methods(http.MethodPost).Subrouter()
	postRouterRestricted.HandleFunc("/upload/post", handler.CreatePost)
	postRouterRestricted.HandleFunc("/upload/document", handler.CreateDocument)
	postRouterRestricted.HandleFunc("/upload/story", handler.CreateStory)
	postRouterRestricted.HandleFunc("/upload/profile-picture", handler.UploadProfilePicture)
	postRouterRestricted.Use(securityMiddleware.Authenticate)

	http.Handle("/", *fs)
}

func main() {
	mediaService := service.NewMediaService()
	mediaHandler := handler.NewMediaHandler(mediaService)
	securityMiddleware := middleware.NewSecurityMiddleware(publicKey)

	sm := mux.NewRouter()
	fs := http.FileServer(http.Dir("./storage"))
	sm.PathPrefix("/storage/").Handler(http.StripPrefix("/storage", fs))

	handleFunc(mediaHandler, securityMiddleware, sm, &fs)

	bindAddress := fmt.Sprintf(":%s", os.Getenv("MEDIA_SERVICE_PORT"))

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
