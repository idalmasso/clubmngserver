package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	model "github.com/idalmasso/clubmngserver/models"
)


var routesUsersProtected =[...]string {"/logout", "/token"}
//addAuthRouterEndpoints add the actual endpoints for auth api
func addUsersRouterEndpoints(r *mux.Router) {
	usersRoute:=r.PathPrefix("/users").Subrouter()
	usersRoute.Use(checkTokenAuthenticationHandler)
	usersRoute.HandleFunc("/", allUsers)
}


func allUsers(w http.ResponseWriter, r *http.Request){
	users,err:=model.GetUsersList(r.Context())
	if err!=nil{
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	sendJSONResponse(w, users, http.StatusOK)
}
