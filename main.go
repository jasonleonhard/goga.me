package main

import (
	"encoding/gob"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/carbocation/interpose"
	gorilla_mux "github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/tlehman/goga.me/dal"
	"github.com/tlehman/goga.me/handlers"
	"github.com/tlehman/goga.me/libenv"
	"github.com/tlehman/goga.me/libunix"
	"github.com/tlehman/goga.me/middlewares"
	"github.com/tylerb/graceful"
	"net/http"
	"time"
)

func init() {
	gob.Register(&dal.UserRow{})
}

// NewApplication is the constructor for Application struct.
func NewApplication() (*Application, error) {
	u, err := libunix.CurrentUser()
	if err != nil {
		return nil, err
	}

	dsn := libenv.EnvWithDefault("DSN", fmt.Sprintf("postgres://%v@localhost:5432/goga.me?sslmode=disable", u))

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, err
	}

	cookieStoreSecret := libenv.EnvWithDefault("COOKIE_SECRET", "Dq7K0yUtpbH7bBrY")

	app := &Application{}
	app.dsn = dsn
	app.db = db
	app.cookieStore = sessions.NewCookieStore([]byte(cookieStoreSecret))

	return app, err
}

// Application is the application object that runs HTTP server.
type Application struct {
	dsn         string
	db          *sqlx.DB
	cookieStore *sessions.CookieStore
}

func (app *Application) middlewareStruct() (*interpose.Middleware, error) {
	middle := interpose.New()
	middle.Use(middlewares.SetDB(app.db))
	middle.Use(middlewares.SetCookieStore(app.cookieStore))

	middle.UseHandler(app.mux())

	return middle, nil
}

func (app *Application) mux() *gorilla_mux.Router {
	MustLogin := middlewares.MustLogin

	router := gorilla_mux.NewRouter()

	router.Handle("/", MustLogin(http.HandlerFunc(handlers.GetHome))).Methods("GET")

	router.HandleFunc("/signup", handlers.GetSignup).Methods("GET")
	router.HandleFunc("/signup", handlers.PostSignup).Methods("POST")
	router.HandleFunc("/login", handlers.GetLogin).Methods("GET")
	router.HandleFunc("/login", handlers.PostLogin).Methods("POST")
	router.HandleFunc("/logout", handlers.GetLogout).Methods("GET")
	router.HandleFunc("/matches", handlers.PostMatches).Methods("POST")
	router.HandleFunc("/matches/{id:[0-9]+}", handlers.GetMatch).Methods("GET")

	router.Handle("/users/{id:[0-9]+}", MustLogin(http.HandlerFunc(handlers.PostPutDeleteUsersID))).Methods("POST", "PUT", "DELETE")

	// Path of static files must be last!
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("static")))

	return router
}

func main() {
	app, err := NewApplication()
	if err != nil {
		logrus.Fatal(err.Error())
	}

	middle, err := app.middlewareStruct()
	if err != nil {
		logrus.Fatal(err.Error())
	}

	serverAddress := libenv.EnvWithDefault("HTTP_ADDR", ":8888")
	certFile := libenv.EnvWithDefault("HTTP_CERT_FILE", "")
	keyFile := libenv.EnvWithDefault("HTTP_KEY_FILE", "")
	drainIntervalString := libenv.EnvWithDefault("HTTP_DRAIN_INTERVAL", "1s")

	drainInterval, err := time.ParseDuration(drainIntervalString)
	if err != nil {
		logrus.Fatal(err.Error())
	}

	srv := &graceful.Server{
		Timeout: drainInterval,
		Server:  &http.Server{Addr: serverAddress, Handler: middle},
	}

	logrus.Infoln("Running HTTP server on " + serverAddress)
	if certFile != "" && keyFile != "" {
		srv.ListenAndServeTLS(certFile, keyFile)
	} else {
		srv.ListenAndServe()
	}
}
