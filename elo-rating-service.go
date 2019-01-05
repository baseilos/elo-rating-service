package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"jozeflang.com/elo-rating-service/dao"
	"net/http"
)

const (
	db_host     = "localhost"
	db_port     = 5432
	db_user     = "postgres"
	db_password = "dev"
	db_dbname   = "elo_service"

	router_port = ":8080"
)

var dbHandler *sql.DB
var err error

func checkDbConnectionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err = dbHandler.Ping()
		if err != nil {
			logrus.Fatalf("Cannot serve %s. No database connection!")
			w.WriteHeader(http.StatusServiceUnavailable)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

// Returns players in the system
func getPlayers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if playerId, ok := vars["id"]; ok {
		//do something here
		logrus.Infof("Getting user with id %s ...", playerId)
		foundPlayers := dao.GetPlayer(dbHandler, playerId)
		if len(foundPlayers) == 0 {
			w.WriteHeader(http.StatusNoContent)
		} else {
			json.NewEncoder(w).Encode(foundPlayers[0])
		}
	} else {
		logrus.Info("Getting all users...")
		json.NewEncoder(w).Encode(dao.GetPlayers(dbHandler))
	}
}

// Returns games played
func getGames(w http.ResponseWriter, r *http.Request) {}

func main() {

	logrus.Info("Starting ELO rating service...")
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		db_host, db_port, db_user, db_password, db_dbname)

	// Connect to database
	logrus.Infof("Connecting to database... %s", psqlInfo)
	dbHandler, err = sql.Open("postgres", psqlInfo)
	defer dbHandler.Close()
	if err != nil {
		panic(err)
	}
	err = dbHandler.Ping()
	if err != nil {
		panic(err)
	}
	logrus.Info("Connection to database successfully established!")

	// Start REST router
	logrus.Infof("Initializing REST router... port=%s", router_port)
	router := mux.NewRouter()
	router.HandleFunc("/players", getPlayers).Methods("GET")
	router.HandleFunc("/player/{id:[0-9]+}", getPlayers).Methods("GET")
	router.HandleFunc("/games", getGames).Methods("GET")
	router.Use(checkDbConnectionMiddleware)
	loggedRouter := handlers.LoggingHandler(logrus.StandardLogger().Out, router)

	logrus.Infof("REST router initialized successfully. Serving %s", router_port)
	err = http.ListenAndServe(router_port, loggedRouter)
	if err != nil {
		logrus.Fatal(err)
	}
}
