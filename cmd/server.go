package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/idalmasso/clubmngserver/app"
	"github.com/idalmasso/clubmngserver/database"
	"github.com/idalmasso/clubmngserver/database/memdb"
	"github.com/idalmasso/clubmngserver/routes"
	"github.com/joho/godotenv"
)

//This will be different (i.e. a map of apps?)
var clubApp app.App
func main(){
	err := godotenv.Load(".env")

  if err != nil {
    log.Fatalf("Error loading .env file")
  }

	var memDB memdb.MemoryDB
	var db database.ClubDb =   &memDB    // Verify that T implements I.
	
	clubApp.InitDB(&db)
	appRoutes:=routes.AppRoutes{App: clubApp}
	
	r := mux.NewRouter()
	r=appRoutes.AddRouteEndpoints(r)
	fs := http.FileServer(http.Dir(os.Getenv("APP_DIR")))
	r.PathPrefix("/").Handler(fs)
	
	http.Handle("/",&corsRouterDecorator{r})
	fmt.Println("Listening")	
	log.Panic(
		http.ListenAndServe(":"+os.Getenv("SERVER_PORT"), nil),
		
	)
}


type corsRouterDecorator struct {
	R *mux.Router
}

func (c *corsRouterDecorator) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if origin := r.Header.Get("Origin"); origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, PATCH")
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")
	}
		// Stop here if its Preflighted OPTIONS request, I just add an OK 
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}


	c.R.ServeHTTP(w, r)
}
