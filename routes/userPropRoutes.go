package routes

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/idalmasso/clubmngserver/common"
	"github.com/idalmasso/clubmngserver/common/userprops"
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


func(appRoutes *AppRoutes) addUserPropsDefinitionRouterEndpoints(r *mux.Router) {
	reqRouter:=r.PathPrefix("/userprop-definitions").Subrouter()
	reqRouter.Use(appRoutes.checkTokenAuthenticationHandler)
	reqRouter.HandleFunc("/", appRoutes.checkUserHasPrivilegeMiddleware(common.SecurityUserPropertyDefinitionView)(appRoutes.getAllUserProps)).Methods("GET")
	reqRouter.HandleFunc("/{userpropname}", appRoutes.checkUserHasPrivilegeMiddleware(common.SecurityUserPropertyDefinitionView)(appRoutes.getSingleUserProp)).Methods("GET")
	reqRouter.HandleFunc("/", appRoutes.checkUserHasPrivilegeMiddleware(common.SecurityUserPropertyDefinitionAdd)(appRoutes.addUserProp)).Methods("POST")
	reqRouter.HandleFunc("/{userpropname}", appRoutes.checkUserHasPrivilegeMiddleware(common.SecurityUserPropertyDefinitionUpdate)(appRoutes.updateUserProp)).Methods("PUT")
	reqRouter.HandleFunc("/{userpropname}", appRoutes.checkUserHasPrivilegeMiddleware(common.SecurityUserPropertyDefinitionRemove)(appRoutes.removeUserProp)).Methods("DELETE")
}

func (appRoutes *AppRoutes)getAllUserProps(w http.ResponseWriter, r *http.Request){
	props,err:=appRoutes.App.GetUserPropertiesList(r.Context())
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
func (appRoutes *AppRoutes)getSingleUserProp(w http.ResponseWriter, r *http.Request){

	params:=mux.Vars(r)
	userPropName, ok:=params["userpropname"]
	if !ok{
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	userProp, err:=appRoutes.App.GetUserProperty(r.Context(), userPropName)
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
func (appRoutes *AppRoutes)addUserProp(w http.ResponseWriter, r *http.Request){
	var userProp UserPropStructValue
	if err:=json.NewDecoder(r.Body).Decode(&userProp); err!=nil{
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	newUserProp, err:=appRoutes.App.AddUserProperty(r.Context(), userProp.Name, userProp.Mandatory, false,userprops.UserPropertyType(userProp.ValueType))
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

func (appRoutes *AppRoutes)updateUserProp(w http.ResponseWriter, r *http.Request){
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
	
	dbProp,  err:=appRoutes.App.UpdateUserProperty(r.Context(), userPropRequest.Name, userPropRequest.Mandatory)
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

func (appRoutes *AppRoutes)removeUserProp(w http.ResponseWriter, r *http.Request){
	
	params:=mux.Vars(r)
	roleName, ok:=params["userpropname"]
	if !ok{
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err:=appRoutes.App.DeleteUserProperty(r.Context(), roleName)
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
