package memdb

import (
	"context"

	"github.com/idalmasso/clubmngserver/common"
	"github.com/idalmasso/clubmngserver/common/userprops"
)

//FindUserPropertyDefinition returns a database.UserPropsDefinition from propertyName
func (db *MemoryDB) FindUserPropertyDefinition(ctx context.Context,propertyName string) (userprops.UserPropertyDefinition, error){
	u, ok:=db.userProps[propertyName]
	if !ok{
		return nil, nil
	}
	return u, nil
}
//GetAllUserPropertyDefinitions returns all users property
func (db *MemoryDB) GetAllUserPropertyDefinitions(ctx context.Context) ([]userprops.UserPropertyDefinition, error){
	var uData []userprops.UserPropertyDefinition
	for _,u:=range(db.userProps){
		uData=append(uData, u)
	}
	return uData, nil
} 

//AddUserPropertyDefinition add a user property to the db
func (db *MemoryDB) AddUserPropertyDefinition(ctx context.Context,userProp userprops.UserPropertyDefinition) (userprops.UserPropertyDefinition,error){
	_, ok:=db.userProps[userProp.GetName()]
	if ok{
		return nil,  common.AlreadyExistsError{ID:userProp.GetName()}
	}
	db.userProps[userProp.GetName()]=userProp
	
	return db.userProps[userProp.GetName()],nil
}
//UpdateUserPropertyDefinition update the user property in the database. In memory does not throw an error
func (db *MemoryDB) UpdateUserPropertyDefinition(ctx context.Context,userProp userprops.UserPropertyDefinition) (userprops.UserPropertyDefinition,error){
	if _, ok:=db.userProps[userProp.GetName()]; !ok{
		return nil, common.NotFoundError{ID:userProp.GetName()}
	}
	db.userProps[userProp.GetName()]=userProp
	return db.userProps[userProp.GetName()], nil
}
//RemoveUserPropertyDefinition remove an actual user property from database
func (db *MemoryDB) RemoveUserPropertyDefinition(ctx context.Context,userPropName string) error{
	delete(db.userProps, userPropName)
	return nil
}
