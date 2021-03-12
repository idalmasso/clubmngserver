package routes

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/idalmasso/clubmngserver/common"
	model "github.com/idalmasso/clubmngserver/models"
)

//AddRouteEndpoints build the mux router with all endpoints
func AddRouteEndpoints(r *mux.Router) *mux.Router {
	apiRouter:=r.PathPrefix("/api").Subrouter()
	addAuthRouterEndpoints(apiRouter)
	addRolesRouterEndpoints(apiRouter)
	addUsersRouterEndpoints(apiRouter)
	return r
}

func sendJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	body, err := json.Marshal(data)
	if err != nil {
		log.Printf("Failed to encode a JSON response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)
	_, err = w.Write(body)
	if err != nil {
		log.Printf("Failed to write the response body: %v", err)
		return
	}
}

func checkTokenAuthenticationHandler(next http.Handler) http.Handler{
	return http.HandlerFunc( func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		bearerToken := strings.Split(header, " ")
		if len(bearerToken)!=2{
			http.Error(w, "Cannot read token", http.StatusBadRequest)
			return
		}
		if bearerToken[0] != "Bearer"{
			http.Error(w, "Error in authorization token. it needs to be in form of 'Bearer <token>'", http.StatusBadRequest)
			return
		}
		
		token, ok :=model.CheckToken(bearerToken[1], true); 
		if !ok{
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if ok && token.Valid {
			
			username, ok := claims["username"].(string)
			if !ok {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return 
			}
			//check if username actually exists
			user,err:= model.FindUser(r.Context(), username);
			if err!=nil{
				var notFoundError common.NotFoundError
				if errors.As(err, &notFoundError){
					w.WriteHeader(http.StatusNotFound)
					w.Write([]byte(notFoundError.Error()))
					return
				}
				http.Error(w, "Unauthorized, error occurred"+err.Error(), http.StatusUnauthorized)
			}
			//Set the username in the request, so I will use it in check after!
			context.Set(r, "user", user)
			context.Set(r, "authenticationtoken", bearerToken[1])
		}
    next.ServeHTTP(w, r)
  })
}


func checkTokenAuthorizationHandler(next http.HandlerFunc) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		bearerToken := strings.Split(header, " ")
		if len(bearerToken)!=2{
			http.Error(w, "Cannot read token", http.StatusBadRequest)
			return
		}
		if bearerToken[0] != "Bearer"{
			http.Error(w, "Error in authorization token. it needs to be in form of 'Bearer <token>'", http.StatusBadRequest)
			return
		}
		
		token, ok :=model.CheckToken(bearerToken[1], false); 
		if !ok{
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if ok && token.Valid {
			
			username, ok := claims["username"].(string)
			if !ok {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return 
			}
			//check if username actually exists
			user,err:= model.FindUser(r.Context(), username); 
			if err!=nil{
				var notFoundError common.NotFoundError
				if errors.As(err, &notFoundError){
					w.WriteHeader(http.StatusNotFound)
					w.Write([]byte(notFoundError.Error()))
					return
				}
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
			}
			//Set the username in the request, so I will use it in check after!
			context.Set(r, "user", user)
			context.Set(r, "authorizationtoken", bearerToken[1])
		}
    next(w, r)
  }
}
func isUsernameContextOk(username string, r *http.Request) bool {
	userCtx, ok:=context.Get(r, "user").(common.UserData)
	if !ok{
		return false
	}
	if userCtx.Username!=username{
		return false
	}
	return true
}
func checkUserHasPrivilegeMiddleware(privileges ...common.SecurityPrivilege) func(http.HandlerFunc) http.HandlerFunc{
	return func (next http.HandlerFunc) http.HandlerFunc{
		return func(w http.ResponseWriter, r *http.Request) {
			userCtx, ok:=context.Get(r, "user").(common.UserData); 
			if !ok{
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			for _,privilege:=range(privileges){
				if authorized, err:=model.IsRoleAuthorized(r.Context(), userCtx.Role, privilege ); err!=nil || !authorized{
					w.WriteHeader(http.StatusUnauthorized)
					if err!=nil{
						w.Write([]byte(err.Error()))
					}
					return
				}
			}
			
			next(w, r)
		}
	}
}

func checkUserHasAtLeastOnePrivilegeMiddleware(privileges ...common.SecurityPrivilege) func(http.HandlerFunc) http.HandlerFunc{
	return func (next http.HandlerFunc) http.HandlerFunc{
		return func(w http.ResponseWriter, r *http.Request) {
			userCtx, ok:=context.Get(r, "user").(common.UserData); 
			if !ok{
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			for _,privilege:=range(privileges){
				if authorized, err:=model.IsRoleAuthorized(r.Context(), userCtx.Role, privilege ); err==nil && authorized{
					next(w,r)
					return
				}
			}
			
			w.WriteHeader(http.StatusUnauthorized)
		}
	}
}
