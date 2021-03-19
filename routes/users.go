package routes

import (
	"errors"
	"net/http"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/idalmasso/clubmngserver/common"
)

//addAuthRouterEndpoints add the actual endpoints for auth api
func (appRoutes *AppRoutes)addUsersRouterEndpoints(r *mux.Router) {
	usersRoute:=r.PathPrefix("/users").Subrouter()
	usersRoute.Use(appRoutes.checkTokenAuthenticationHandler)
	usersRoute.HandleFunc("/", appRoutes.checkUserHasPrivilegeMiddleware(common.SecurityAllUsersView)(appRoutes.getAllUsers)).Methods("GET")
	usersRoute.HandleFunc("/{userName}", appRoutes.checkUserHasPrivilegeMiddleware(common.SecurityAllUsersView, common.SecurityLinkedUsersView, common.SecuritySelfUserView)(appRoutes.getAllUsers)).Methods("GET")
	//usersRoute.HandleFunc("/", allUsers)
}


func (appRoutes *AppRoutes)getAllUsers(w http.ResponseWriter, r *http.Request){
	users,err:=appRoutes.App.GetUsersList(r.Context())
	if err!=nil{
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	sendJSONResponse(w, users, http.StatusOK)
}

func (appRoutes *AppRoutes)getSingleUser(w http.ResponseWriter, r *http.Request){
	params:=mux.Vars(r)
	userName, ok:=params["userName"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	canSeeAllUsers,err:=appRoutes.App.IsRoleAuthorized(r.Context(), context.Get(r, "user").(common.UserData).Role, common.SecurityAllUsersView)
	if err!=nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	canSeeLinkedUsers,err:=appRoutes.App.IsRoleAuthorized(r.Context(), context.Get(r, "user").(common.UserData).Role, common.SecurityLinkedUsersView)
	if err!=nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	
	if context.Get(r, "user").(common.UserData).Username!=userName && !canSeeAllUsers &&	!canSeeLinkedUsers{
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	user, err:=appRoutes.App.FindUser(r.Context(),userName)
	if err!=nil {
		var notFoundError common.NotFoundError
		if errors.As(err, &notFoundError){
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(notFoundError.Error()))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	sendJSONResponse(w, user, http.StatusOK)
}
