package routes

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/idalmasso/clubmngserver/common"
	"github.com/idalmasso/clubmngserver/common/userprops"
	model "github.com/idalmasso/clubmngserver/models"
)
type UserPropStructValue struct{
	Name string `json:"name"`
	ValueType string `json:"type"`
	Mandatory bool `json:"mandatory"`
}

func (userProperty *UserPropStructValue) initFromDBUserProp(property userprops.UserPropertyDefinition){
	userProperty.Name=property.GetName()
	userProperty.Mandatory=property.IsMandatory()
	userProperty.ValueType=string(property.GetTypeString())
}


func addUserPropsDefinitionRouterEndpoints(r *mux.Router) {
	reqRouter:=r.PathPrefix("/userprop-definitions").Subrouter()
	reqRouter.Use(checkTokenAuthenticationHandler)
	reqRouter.HandleFunc("/", checkUserHasPrivilegeMiddleware(common.SecurityUserPropertyDefinitionView)(getAllUserProps)).Methods("GET")
	reqRouter.HandleFunc("/{userpropname}", checkUserHasPrivilegeMiddleware(common.SecurityUserPropertyDefinitionView)(getSingleUserProp)).Methods("GET")
	reqRouter.HandleFunc("/", checkUserHasPrivilegeMiddleware(common.SecurityUserPropertyDefinitionAdd)(addUserProp)).Methods("POST")
	reqRouter.HandleFunc("/{userpropname}", checkUserHasPrivilegeMiddleware(common.SecurityUserPropertyDefinitionUpdate)(updateUserProp)).Methods("PUT")
	reqRouter.HandleFunc("/{userpropname}", checkUserHasPrivilegeMiddleware(common.SecurityUserPropertyDefinitionRemove)(removeUserProp)).Methods("DELETE")
}

func getAllUserProps(w http.ResponseWriter, r *http.Request){
	props,err:=model.GetUserPropertiesList(r.Context())
	if err!=nil{
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	var propsReturn [] UserPropStructValue
	for _,prop:=range(props){
		var newProp UserPropStructValue 
		newProp.initFromDBUserProp(prop)
		propsReturn=append(propsReturn,newProp )
	}
	sendJSONResponse(w, propsReturn, http.StatusOK)
	
}
func getSingleUserProp(w http.ResponseWriter, r *http.Request){

	params:=mux.Vars(r)
	userPropName, ok:=params["userpropname"]
	if !ok{
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	userProp, err:=model.GetUserProperty(r.Context(), userPropName)
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
	var userPropRequest UserPropStructValue
	userPropRequest.initFromDBUserProp(userProp)
	sendJSONResponse(w, userPropRequest, http.StatusOK)
}
func addUserProp(w http.ResponseWriter, r *http.Request){
	var userProp UserPropStructValue
	if err:=json.NewDecoder(r.Body).Decode(&userProp); err!=nil{
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	newUserProp, err:=model.AddUserProperty(r.Context(), userProp.Name, userProp.Mandatory, userprops.UserPropertyType(userProp.ValueType))
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
	userProp.initFromDBUserProp(newUserProp)
	sendJSONResponse(w, userProp, http.StatusOK)
}

func updateUserProp(w http.ResponseWriter, r *http.Request){
	var userPropRequest UserPropStructValue
	if err:=json.NewDecoder(r.Body).Decode(&userPropRequest); err!=nil{
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	params:=mux.Vars(r)
	userPropName, ok:=params["userpropname"]
	if !ok || userPropName!=userPropRequest.Name{
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	
	dbProp,  err:=model.UpdateUserProperty(r.Context(), userPropRequest.Name, userPropRequest.Mandatory)
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
	
	userPropRequest.initFromDBUserProp(dbProp)
	sendJSONResponse(w, userPropRequest, http.StatusOK)
}

func removeUserProp(w http.ResponseWriter, r *http.Request){
	
	params:=mux.Vars(r)
	roleName, ok:=params["userpropname"]
	if !ok{
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err:=model.DeleteUserProperty(r.Context(), roleName)
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
