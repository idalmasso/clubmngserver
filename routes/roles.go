package routes

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)


var routesRolesProtected =[...]string {"/logout", "/token"}
//addAuthRouterEndpoints add the actual endpoints for auth api
func addRolesRouterEndpoints(r *mux.Router) {
	authRouter:=r.PathPrefix("/roles").Subrouter()
	
	reqRouter:=authRouter.MatcherFunc(func(r *http.Request,  rm *mux.RouteMatch) bool{
		for _,route:=range(routesAuthProtected){
			if strings.Contains(r.RequestURI, route){
				return true
			}
		}
		return false
	}).Subrouter()
	
	
	reqRouter.Use(checkTokenAuthenticationHandler)

}


