package models

import (
	"context"

	"github.com/idalmasso/clubmngserver/common"
	"github.com/idalmasso/clubmngserver/common/userprops"
)

//AddUserProperty tries to add a user property to the database
func AddUserProperty(ctx context.Context, userPropertyName string, mandatory bool, isSystem bool, userPropertyType userprops.UserPropertyType) (userprops.UserPropertyDefinition,error) {
	prop, err:=db.FindUserPropertyDefinition(ctx, userPropertyName)
	if err!=nil{
		return nil,err
	}
	if prop!=nil{
		return nil,  common.AlreadyExistsError{ID:userPropertyName}
	}
	newProp, err:= userprops.NewUserPropertyDefinition(userPropertyType)
	newProp.SetName(userPropertyName)
	newProp.SetMandatory(mandatory)
	newProp.SetIsSystem(isSystem)
	newProp,err=db.AddUserPropertyDefinition(ctx,newProp)
	if err!=nil{
		return nil,err
	}
	return newProp, nil
}
//DeleteUserProperty deletes a user property (not its values!). Cannot delete system properties
func DeleteUserProperty(ctx context.Context, userPropertyName string) error {
	property, err:=db.FindUserPropertyDefinition(ctx, userPropertyName)
	if err!=nil{
		return err
	}
	if property==nil{
		return common.NotFoundError{ID:userPropertyName}
	}
	if property.IsSystem(){
		return common.BadRequestParametersError{Message: "Cannot delete system property"}
	}
	return db.RemoveRole(ctx, userPropertyName)
}

//UpdateUserProperty updates userProperty definitions. Cannot update or delete system properties
func UpdateUserProperty(ctx context.Context, userPropertyDefinitionName string, isMandatory bool ) (userprops.UserPropertyDefinition, error) {
	prop, err:=db.FindUserPropertyDefinition(ctx, userPropertyDefinitionName)
	if err!=nil{
		return nil,err
	}
	if prop==nil{
		return nil,common.NotFoundError{ID:userPropertyDefinitionName}
	}
	prop.SetMandatory(isMandatory)
	
	_, err=db.UpdateUserPropertyDefinition(ctx, prop)
	
	return prop,err
}
//GetUserPropertiesList return a list of all propert definitions
func GetUserPropertiesList(ctx context.Context) ([]userprops.UserPropertyDefinition, error){
	props, err:=db.GetAllUserPropertyDefinitions(ctx)
	return props,err
}

//GetUserProperty finds the property definition and return it if exists
func GetUserProperty(ctx context.Context, userPropertyName string) (userprops.UserPropertyDefinition, error){
	property, err:=db.FindUserPropertyDefinition(ctx, userPropertyName)
	if property==nil{
		return nil, common.NotFoundError{ID: userPropertyName}
	}
	return property,err
}
