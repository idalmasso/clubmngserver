package routes

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	model "github.com/idalmasso/clubmngserver/models"
)
type UserPassword struct{
	Username string `json:'username'`
	Password string `json:'password'`
}
type TokensResponse struct{
	AuthenticationToken string `json:'authentication'`
	AuthorizationToken string `json:'authorization'`
}
//addAuthRouterEndpoints add the actual endpoints for auth api
func addAuthRouterEndpoints(r *mux.Router) *mux.Router {
	r.HandleFunc("/api/auth/login", login).Methods("POST")
	r.HandleFunc("/api/auth/logout", checkTokenAuthenticationHandler(logout)).Methods("POST")
	r.HandleFunc("/api/auth/create-user", createUser).Methods("POST")
	r.HandleFunc("/api/auth/token", checkTokenAuthenticationHandler(getTokenByToken)).Methods("GET")
	return r
}


func login(w http.ResponseWriter, r *http.Request){
	var u UserPassword 
	err:=json.NewDecoder(r.Body).Decode(&u)
	if err!=nil{
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user:=model.FindUser(r.Context(), u.Username)
	if user==nil{
		w.WriteHeader(http.StatusNotFound)
		return
	}
	authent, author, err:=user.TryAuthenticate(r.Context(), u.Password)
	if err!=nil{
		w.WriteHeader(http.StatusForbidden)
		return
	}
	t:=TokensResponse{AuthenticationToken: authent, AuthorizationToken: author }
	sendJSONResponse(w, t, http.StatusAccepted)
}

func logout(w http.ResponseWriter, r *http.Request){
	var u *model.UserDetails
	u =context.Get(r, "user").(*model.UserDetails)
	var authenticationToken string
	authenticationToken = context.Get(r, "authenticationtoken").(string)
	u.RemoveUserAuthentication(authenticationToken)
}

func createUser(w http.ResponseWriter, r *http.Request){

}


func getTokenByToken(w http.ResponseWriter, r *http.Request){

}


