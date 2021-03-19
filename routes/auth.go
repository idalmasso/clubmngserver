package routes

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/idalmasso/clubmngserver/common"
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


//addAuthRouterEndpoints add the actual endpoints for auth api
func (appRoutes *AppRoutes)addAuthRouterEndpoints(r *mux.Router) {
	authRouter:=r.PathPrefix("/auth").Subrouter()
	var routesAuthProtected =[...]string {"/logout", "/token"}
	reqRouter:=authRouter.MatcherFunc(func(r *http.Request,  rm *mux.RouteMatch) bool{
		for _,route:=range(routesAuthProtected){
			if strings.Contains(r.RequestURI, route){
				return true
			}
		}
		return false
	}).Subrouter()
	
	authRouter.HandleFunc("/login", appRoutes.login).Methods("POST")
	authRouter.HandleFunc("/signin", appRoutes.createUser).Methods("POST")
	reqRouter.Use(appRoutes.checkTokenAuthenticationHandler)
	reqRouter.HandleFunc("/logout", appRoutes.logout).Methods("POST")
	reqRouter.HandleFunc("/token", appRoutes.getTokenByToken).Methods("GET")
}


func (appRoutes *AppRoutes)login(w http.ResponseWriter, r *http.Request){
	var u UserPassword 
	err:=json.NewDecoder(r.Body).Decode(&u)
	if err!=nil{
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user,err:=appRoutes.App.FindUser(r.Context(), u.Username)
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
	authent, author, err:=appRoutes.App.TryAuthenticate(r.Context(),*user, u.Password)
	if err!=nil{
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	t:=TokensResponse{AuthenticationToken: authent, AuthorizationToken: author }
	sendJSONResponse(w, t, http.StatusAccepted)
}

func (appRoutes *AppRoutes)logout(w http.ResponseWriter, r *http.Request){
	u :=context.Get(r, "user").(common.UserData)
	var authenticationToken string
	authenticationToken = context.Get(r, "authenticationtoken").(string)
	err:=appRoutes.App.RemoveUserAuthentication(r.Context(), u.Username,authenticationToken)
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

func (appRoutes *AppRoutes)createUser(w http.ResponseWriter, r *http.Request){
	var u UserPassword 
	err:=json.NewDecoder(r.Body).Decode(&u)
	if err!=nil{
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var newUser *common.UserData
	newUser.Username=u.Username

	newUser, err=appRoutes.App.AddUser(r.Context(),*newUser, u.Password)
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


func (appRoutes *AppRoutes)getTokenByToken(w http.ResponseWriter, r *http.Request){

}


