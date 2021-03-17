package models

import (
	"context"

	"github.com/idalmasso/clubmngserver/common"
	"github.com/idalmasso/clubmngserver/database"
)

type passwordRole struct{
	password string
	role string
}

var db database.ClubDb
var rolesToBeAdded=map[string][]common.SecurityPrivilege{"Admin": {common.SecurityAdmin}}
var usersToBeAdded=map[string]passwordRole{"admin":{password: "Abcd1234", role:"Admin"}}


//InitDB calls the actual initialization of the database and also insert the basic data in it
func InitDB(database *database.ClubDb){
	db=*database
	db.Init()
	addDatabaseInitialSetup()
}

func addDatabaseInitialSetup(){ 
	initTryAddRoles()
	initTryAddUsers()
	initTryAddStandardUserFields()
}
func initTryAddRoles(){
	for adminRoleName, privileges:=range(rolesToBeAdded){
		
		r, err:=db.FindRole(context.Background(), adminRoleName)
		if err!=nil{
			panic("Cannot find role admin on initializations because of error "+ err.Error())
		}
		if r==nil{
			var adminRole common.SecurityRole
			adminRole.Name=adminRoleName
			adminRole.Privileges= privileges
			db.AddRole(context.Background(), adminRole) 
		}
		
	}
}

func initTryAddUsers(){
	for adminUserName, rolePass:=range(usersToBeAdded){
		uData, err:=db.FindUser(context.Background(), adminUserName)
		if err!=nil{
			panic("Cannot find user" +adminUserName+" on initializations because of error "+ err.Error())
		}
		if uData==nil{
			uData = &common.UserData{}
			uData.Username=adminUserName
			uData, err=db.AddUser(context.Background(), *uData)
			if err!=nil{
				panic("Cannot create "+ adminUserName+" user")
			}
			
			if err=uData.SetPassword(context.Background(),rolePass.password); err!=nil{
				panic("Cannot Set user "+adminUserName+" password on initializations because of error "+ err.Error())
			}
			role, err:=db.FindRole(context.Background(), rolePass.role)
			if err!=nil {
				panic("Cannot find again role "+ rolePass.role +" "+err.Error())
			}
			if role==nil{
				panic("Cannot find again role "+ rolePass.role)
			}
			uData.AddRole(context.Background(), *role)
			if _, err=db.UpdateUser(context.Background(), *uData); err!=nil{
				panic("Cannot update user admin on initializations because of error "+ err.Error())
			}
		}
	}
}

func initTryAddStandardUserFields(){
	return
}
