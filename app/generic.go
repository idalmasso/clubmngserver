package app

import (
	"context"
	"errors"

	"github.com/idalmasso/clubmngserver/common"
	"github.com/idalmasso/clubmngserver/common/userprops"
	"github.com/idalmasso/clubmngserver/database"
)



type passwordRole struct{
	password string
	role string
}
//App is the actual app, contains the application database and the logic. 
//
//This will be used to have multiple instances of this running together
type App struct{
	Name string
	db database.ClubDb
}
//TODO: Move this code somewhere else
var rolesToBeAdded=map[string][]common.SecurityPrivilege{"Admin": {common.SecurityAdmin}}
var usersToBeAdded=map[string]passwordRole{"admin":{password: "Abcd1234", role:"Admin"}}
func buildUserPropertiesToBeAdded() map[string]userprops.UserPropertyDefinition {
	userPropertiesToBeAdded := make(map[string]userprops.UserPropertyDefinition)
	addToMapProperty("email", userprops.UserTypeString, true, true, userPropertiesToBeAdded)
	addToMapProperty("age", userprops.UserTypeInt64, false, true, userPropertiesToBeAdded)
	addToMapProperty("emailVerified", userprops.UserTypeBool, false, true, userPropertiesToBeAdded)
	return userPropertiesToBeAdded
}
func addToMapProperty(name string, propType userprops.UserPropertyType, isMandatory bool, isSystem bool, userProperties map[string]userprops.UserPropertyDefinition){
	p, _:=userprops.NewUserPropertyDefinition(userprops.UserTypeString)
	p.SetIsSystem(isSystem)
	p.SetMandatory(isMandatory)
	p.SetName(name)
	userProperties[name]=p
}
//InitDB calls the actual initialization of the database and also insert the basic data in it
func(app *App)InitDB(database *database.ClubDb){
	app.db=*database
	app.db.Init()
	app.addDatabaseInitialSetup()
}

func (app *App)addDatabaseInitialSetup(){ 
	app.initTryAddRoles()
	app.initTryAddUsers()
	app.initTryAddStandardUserFields()
}
func (app *App)initTryAddRoles(){
	for adminRoleName, privileges:=range(rolesToBeAdded){
		
		r, err:=app.db.FindRole(context.Background(), adminRoleName)
		if err!=nil{
			panic("Cannot find role admin on initializations because of error "+ err.Error())
		}
		if r==nil{
			var adminRole common.SecurityRole
			adminRole.Name=adminRoleName
			adminRole.Privileges= privileges
			app.db.AddRole(context.Background(), adminRole) 
		}
		
	}
}

func (app *App)initTryAddUsers(){
	for adminUserName, rolePass:=range(usersToBeAdded){
		uData, err:=app.db.FindUser(context.Background(), adminUserName)
		if err!=nil{
			panic("Cannot find user" +adminUserName+" on initializations because of error "+ err.Error())
		}
		if uData==nil{
			uData = &common.UserData{}
			uData.Username=adminUserName
			uData, err=app.db.AddUser(context.Background(), *uData)
			if err!=nil{
				panic("Cannot create "+ adminUserName+" user")
			}
			
			if err=uData.SetPassword(context.Background(),rolePass.password); err!=nil{
				panic("Cannot Set user "+adminUserName+" password on initializations because of error "+ err.Error())
			}
			role, err:=app.db.FindRole(context.Background(), rolePass.role)
			if err!=nil {
				panic("Cannot find again role "+ rolePass.role +" "+err.Error())
			}
			if role==nil{
				panic("Cannot find again role "+ rolePass.role)
			}
			uData.AddRole(context.Background(), *role)
			if _, err=app.db.UpdateUser(context.Background(), *uData); err!=nil{
				panic("Cannot update user admin on initializations because of error "+ err.Error())
			}
		}
	}
}

func(app *App) initTryAddStandardUserFields(){
	userPropertiesToBeAdded := buildUserPropertiesToBeAdded()
	for _, field:=range(userPropertiesToBeAdded){
		if _, err :=app.db.AddUserPropertyDefinition(context.Background(), field); err!=nil{
			if !errors.As(err, &common.AlreadyExistsError{}){
				panic("Cannot insert user prop: "+err.Error())
			}
		}
	}
	return
}
