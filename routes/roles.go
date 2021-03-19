package routes

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/idalmasso/clubmngserver/common"
)

//RoleStructValue is the value sent and received from web service
type RoleStructValue struct{
	Name string `json:"name"`
	Privileges []string `json:"privileges"` 
}

func (roleValue *RoleStructValue) initFromDBRole(role common.SecurityRole){
	roleValue.Name=role.Name
	for _,privilege:=range(role.Privileges){
		pStr, err:=privilege.String()
		if err==nil{
			roleValue.Privileges=append(roleValue.Privileges, pStr)
		}
	}
}

func  (roleValue *RoleStructValue)  toDBRole()( common.SecurityRole) {
	role:=common.SecurityRole{Name: roleValue.Name}
	for _, privilegeStr:=range(roleValue.Privileges){
		privilege, err:=common.StringToSecurityPrivilege(privilegeStr)
		if err==nil{
			role.Privileges=append(role.Privileges, privilege)
		}
	}
	return role
}

func(appRoutes *AppRoutes) addRolesRouterEndpoints(r *mux.Router) {
	reqRouter:=r.PathPrefix("/roles").Subrouter()
	privRouter:=r.PathPrefix("/privileges").Subrouter()
	reqRouter.Use(appRoutes.checkTokenAuthenticationHandler)
	reqRouter.HandleFunc("/", appRoutes.checkUserHasPrivilegeMiddleware(common.SecurityRolesView)(appRoutes.allRoles)).Methods("GET")
	reqRouter.HandleFunc("/{roleName}", appRoutes.checkUserHasPrivilegeMiddleware(common.SecurityRolesView)(appRoutes.viewRole)).Methods("GET")
	reqRouter.HandleFunc("/", appRoutes.checkUserHasPrivilegeMiddleware(common.SecurityRolesAdd)(appRoutes.addRole)).Methods("POST")
	reqRouter.HandleFunc("/{roleName}", appRoutes.checkUserHasPrivilegeMiddleware(common.SecurityRolesUpdate)(appRoutes.updateRole)).Methods("PUT")
	reqRouter.HandleFunc("/{roleName}", appRoutes.checkUserHasPrivilegeMiddleware(common.SecurityRolesDelete)(appRoutes.removeRole)).Methods("DELETE")
	privRouter.HandleFunc("/", appRoutes.checkUserHasPrivilegeMiddleware(common.SecurityRolesView)(listPrivileges)).Methods("GET")
}

func listPrivileges(w http.ResponseWriter, r *http.Request){
	privileges:=common.ListPrivileges()
	sendJSONResponse(w,privileges,http.StatusOK )
}

func(appRoutes *AppRoutes) allRoles(w http.ResponseWriter, r *http.Request){
	roles,err:=appRoutes.App.GetAllRoles(r.Context())
	if err!=nil{
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	var rolesReturn [] RoleStructValue
	for _,role:=range(roles){
		var newRole RoleStructValue 
		newRole.initFromDBRole(role)
		rolesReturn=append(rolesReturn,newRole )
	}
	sendJSONResponse(w, rolesReturn, http.StatusOK)
}
func(appRoutes *AppRoutes) viewRole(w http.ResponseWriter, r *http.Request){
	params:=mux.Vars(r)
	roleName, ok:=params["roleName"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	role, err:=appRoutes.App.GetRole(r.Context(),roleName)
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
	var roleRequest RoleStructValue
	roleRequest.initFromDBRole(*role)
	sendJSONResponse(w, roleRequest, http.StatusOK)
}
func (appRoutes *AppRoutes)addRole(w http.ResponseWriter, r *http.Request){
	var roleRequest RoleStructValue
	if err:=json.NewDecoder(r.Body).Decode(&roleRequest); err!=nil{
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	newRole := roleRequest.toDBRole()
	addedUser, err:=appRoutes.App.AddRole(r.Context(), newRole.Name, newRole.Privileges...)
	if err!=nil{
		var alreadyExists common.AlreadyExistsError
		if errors.As(err, &alreadyExists){
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte(alreadyExists.Error()))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	roleRequest.initFromDBRole(*addedUser)
	sendJSONResponse(w, roleRequest, http.StatusOK)
}
func (appRoutes *AppRoutes)updateRole(w http.ResponseWriter, r *http.Request){
	var roleRequest RoleStructValue
	if err:=json.NewDecoder(r.Body).Decode(&roleRequest); err!=nil{
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	params:=mux.Vars(r)
	roleName, ok:=params["roleName"]
	if !ok || roleName!=roleRequest.Name{
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	newRole:=roleRequest.toDBRole()
	_,err:=appRoutes.App.UpdateRole(r.Context(), newRole.Name, newRole.Privileges...)
	if err!=nil{
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
	
	roleRequest.initFromDBRole(newRole)
	sendJSONResponse(w, roleRequest, http.StatusOK)
}
func(appRoutes *AppRoutes) removeRole(w http.ResponseWriter, r *http.Request){
	
	params:=mux.Vars(r)
	roleName, ok:=params["roleName"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err:=appRoutes.App.DeleteRole(r.Context(), roleName)
	if err!=nil{
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
	w.WriteHeader(http.StatusOK)
}
