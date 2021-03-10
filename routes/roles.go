package routes

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/idalmasso/clubmngserver/common"
	model "github.com/idalmasso/clubmngserver/models"
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

func addRolesRouterEndpoints(r *mux.Router) {
	reqRouter:=r.PathPrefix("/roles").Subrouter()
	privRouter:=r.PathPrefix("/privileges").Subrouter()
	reqRouter.Use(checkTokenAuthenticationHandler)
	reqRouter.HandleFunc("/", checkUserHasPrivilegeMiddleware(common.SecurityRolesView)(allRoles)).Methods("GET")
	reqRouter.HandleFunc("/{roleName}", checkUserHasPrivilegeMiddleware(common.SecurityRolesView)(viewRole)).Methods("GET")
	reqRouter.HandleFunc("/", checkUserHasPrivilegeMiddleware(common.SecurityRolesUpdate)(addRole)).Methods("POST")
	reqRouter.HandleFunc("/{roleName}", checkUserHasPrivilegeMiddleware(common.SecurityRolesUpdate)(updateRole)).Methods("PUT")
	reqRouter.HandleFunc("/{roleName}", checkUserHasPrivilegeMiddleware(common.SecurityRolesUpdate)(removeRole)).Methods("DELETE")
	privRouter.HandleFunc("/", checkUserHasPrivilegeMiddleware(common.SecurityRolesView)(listPrivileges)).Methods("GET")
}

func listPrivileges(w http.ResponseWriter, r *http.Request){
	privileges:=common.ListPrivileges()
	sendJSONResponse(w,privileges,http.StatusOK )
}

func allRoles(w http.ResponseWriter, r *http.Request){
	roles,err:=model.GetAllRoles(r.Context())
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
func viewRole(w http.ResponseWriter, r *http.Request){
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
	role, err:=model.GetRole(r.Context(), roleRequest.Name)
	if err!=nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	if role==nil{
		w.WriteHeader(http.StatusNotFound)
		return
	}
	roleRequest.initFromDBRole(*role)
	sendJSONResponse(w, roleRequest, http.StatusOK)
}
func addRole(w http.ResponseWriter, r *http.Request){
	var roleRequest RoleStructValue
	if err:=json.NewDecoder(r.Body).Decode(&roleRequest); err!=nil{
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	newRole := roleRequest.toDBRole()
	addedUser, err:=model.AddRole(r.Context(), newRole.Name, newRole.Privileges...)
	if err!=nil{
		if strings.Contains(err.Error(), "already exists"){
			roleRequest.initFromDBRole(*addedUser)
			sendJSONResponse(w,roleRequest, http.StatusSeeOther)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	roleRequest.initFromDBRole(*addedUser)
	sendJSONResponse(w, roleRequest, http.StatusOK)
}
func updateRole(w http.ResponseWriter, r *http.Request){
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
	err:=model.UpdateRole(r.Context(), newRole.Name, newRole.Privileges...)
	if err!=nil{
		if strings.Contains(err.Error(),  "does not exists"){
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	
	roleRequest.initFromDBRole(newRole)
	sendJSONResponse(w, roleRequest, http.StatusOK)
}
func removeRole(w http.ResponseWriter, r *http.Request){
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
	err:=model.DeleteRole(r.Context(), roleRequest.Name)
	if err!=nil{
		if strings.Contains(err.Error(),  "does not exists"){
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
}
