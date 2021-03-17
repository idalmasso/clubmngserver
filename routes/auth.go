package routes

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/idalmasso/clubmngserver/common"
	model "github.com/idalmasso/clubmngserver/models"
)

//UserPassword is the type the user send with username and password for authentication
type UserPassword struct{
	Username string `json:"username"`
	Password string `json:"password"`
}
//TokensResponse is the response type with one or two authentication tokens
type TokensResponse struct{
	AuthenticationToken string `json:"authentication"`
	AuthorizationToken string `json:"authorization"`
}

var routesAuthProtected =[...]string {"/logout", "/token"}
//addAuthRouterEndpoints add the actual endpoints for auth api
func addAuthRouterEndpoints(r *mux.Router) {
	authRouter:=r.PathPrefix("/auth").Subrouter()
	
	reqRouter:=authRouter.MatcherFunc(func(r *http.Request,  rm *mux.RouteMatch) bool{
		for _,route:=range(routesAuthProtected){
			if strings.Contains(r.RequestURI, route){
				return true
			}
		}
		return false
	}).Subrouter()
	
	authRouter.HandleFunc("/login", login).Methods("POST")
	authRouter.HandleFunc("/signin", createUser).Methods("POST")
	reqRouter.Use(checkTokenAuthenticationHandler)
	reqRouter.HandleFunc("/logout", logout).Methods("POST")
	reqRouter.HandleFunc("/token", getTokenByToken).Methods("GET")
}


func login(w http.ResponseWriter, r *http.Request){
	var u UserPassword 
	err:=json.NewDecoder(r.Body).Decode(&u)
	if err!=nil{
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user,err:=model.FindUser(r.Context(), u.Username)
	if err!=nil{
		var notFoundError common.NotFoundError
		if errors.As(err, &notFoundError){
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(notFoundError.Error()))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	authent, author, err:=model.TryAuthenticate(r.Context(),*user, u.Password)
	if err!=nil{
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	t:=TokensResponse{AuthenticationToken: authent, AuthorizationToken: author }
	sendJSONResponse(w, t, http.StatusAccepted)
}

func logout(w http.ResponseWriter, r *http.Request){
	u :=context.Get(r, "user").(common.UserData)
	var authenticationToken string
	authenticationToken = context.Get(r, "authenticationtoken").(string)
	err:=model.RemoveUserAuthentication(r.Context(), u.Username,authenticationToken)
	if err!=nil{
		if errors.As(err,&common.NotFoundError{}){
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func createUser(w http.ResponseWriter, r *http.Request){
	var u UserPassword 
	err:=json.NewDecoder(r.Body).Decode(&u)
	if err!=nil{
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var newUser *common.UserData
	newUser.Username=u.Username

	newUser, err=model.AddUser(r.Context(),*newUser, u.Password)
	if err!=nil{
		var alreadyExists common.AlreadyExistsError
		if errors.As(err, &alreadyExists){
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte(alreadyExists.Error()))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	sendJSONResponse(w, newUser, http.StatusCreated)
}


func getTokenByToken(w http.ResponseWriter, r *http.Request){

}


