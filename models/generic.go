package models

import (
	"context"

	"github.com/idalmasso/clubmngserver/common"
	"github.com/idalmasso/clubmngserver/database"
)


var db database.ClubDb
//InitDB calls the actual initialization of the database and also insert the basic data in it
func InitDB(database *database.ClubDb){
	db=*database
	db.Init()
	addDatabaseInitialSetup()
}

func addDatabaseInitialSetup(){ 
	adminRoleName:="Admin"
	adminUserName:="administrator"
	adminUserPassword:="Abcd1234"
	initTryAddAdminRole(adminRoleName)
	initTryAddAdmin(adminRoleName, adminUserName, adminUserPassword)
	initTryAddStandardUSerFields()
}
func initTryAddAdminRole(adminRoleName string){
	r, err:=db.FindRole(context.Background(), adminRoleName)
	if err!=nil{
		panic("Cannot find role admin on initializations because of error "+ err.Error())
	}
	if r==nil{
		var adminRole common.SecurityRole
		adminRole.Name=adminRoleName
		adminRole.Privileges= make([]common.SecurityPrivilege, 0)
		adminRole.Privileges = append(adminRole.Privileges, common.SecurityAdmin)
		db.AddRole(context.Background(), adminRole) 
	}
}

func initTryAddAdmin(adminRoleName, adminUserName, adminUserPassword string){
	uData, err:=db.FindUser(context.Background(), adminUserName)
	if err!=nil{
		panic("Cannot find user admin on initializations because of error "+ err.Error())
	}
	if uData==nil{
		uData.Username=adminUserName
		uData, err=db.AddUser(context.Background(), *uData)
		if err!=nil{
			panic("Cannot create admin user")
		}
		
		if err=uData.SetPassword(context.Background(),adminUserPassword); err!=nil{
			panic("Cannot Set user admin password on initializations because of error "+ err.Error())
		}
		role, err:=db.FindRole(context.Background(), adminRoleName)
		if err!=nil{
			panic("Cannot find again role admin "+err.Error())
		}

		uData.AddRole(context.Background(), *role)
		if _, err=db.UpdateUser(context.Background(), *uData); err!=nil{
			panic("Cannot update user admin on initializations because of error "+ err.Error())
		}
	}
}

func initTryAddStandardUSerFields(){
	return
}
