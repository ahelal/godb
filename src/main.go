package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/jinzhu/configor"
)

var Config = struct {
	DB struct {
		File   string `default:"db.json"`
		EncKey string
		Sync   int
	}
	Web struct {
		Listen       string `default:"0.0.0.0"`
		Port         uint   `default:"10000"`
		BasicAuth    bool
		User         string
		Password     string
		GracefulWait time.Duration `default:"30s"`
	}
}{}

func protectedOrNot(h http.HandlerFunc) http.HandlerFunc {
	if Config.Web.BasicAuth {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			s := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
			if len(s) != 2 {
				http.Error(w, "Not authorized", 401)
				return
			}
			b, err := base64.StdEncoding.DecodeString(s[1])
			if err != nil {
				http.Error(w, err.Error(), 401)
				return
			}
			pair := strings.SplitN(string(b), ":", 2)
			if len(pair) != 2 {
				http.Error(w, "Not authorized", 401)
				return
			}
			if pair[0] != Config.Web.User || pair[1] != Config.Web.Password {
				http.Error(w, "Not authorized", 401)
				return
			}
			h.ServeHTTP(w, r)
		}
	}
	return func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
	}
}
func use(h http.HandlerFunc, middleware ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	for _, m := range middleware {
		h = m(h)
	}
	return h
}

func myHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Authenticated!"))
	return
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
}

func main() {
	var configFile string
	if len(os.Args) > 1 {
		configFile = os.Args[1]
	} else {
		configFile = "config.yml"
	}
	err := configor.Load(&Config, configFile)
	if err != nil {
		panic(err)
	}
	initDB()
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	// kV
	myRouter.HandleFunc("/k/keys", use(listKey, protectedOrNot))
	myRouter.HandleFunc("/k/inc/{key}", use(incByKey, protectedOrNot)).Methods("POST", "GET", "DELETE")
	myRouter.HandleFunc("/k/key/{key}", use(keyValue, protectedOrNot)).Methods("POST", "GET", "DELETE")
	// queue
	myRouter.HandleFunc("/q/keys", use(listQueues, protectedOrNot))
	myRouter.HandleFunc("/q/len/{key}", use(lenQ, protectedOrNot))
	myRouter.HandleFunc("/q/push/{key}", use(pushByKey, protectedOrNot)).Methods("POST")
	myRouter.HandleFunc("/q/pop/{key}", use(popByKey, protectedOrNot)).Methods("POST", "GET")

	connectString := fmt.Sprintf("%s:%d", Config.Web.Listen, Config.Web.Port)
	fmt.Println("Connection string", connectString)
	srv := &http.Server{
		Addr:         connectString,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      myRouter,
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), Config.Web.GracefulWait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.

	ticker.Stop()
	dbSync(true)
	log.Println("shutting down")
	os.Exit(0)
}
