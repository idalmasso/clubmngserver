package models

import (
	"context"
	"fmt"

	"github.com/idalmasso/clubmngserver/database"
)

func AddRole(ctx context.Context, roleName string, privileges ...database.SecurityPrivilege) ( *database.SecurityRole,error) {
	role, err:=db.FindRole(ctx, roleName)
	if err!=nil{
		return nil,err
	}
	if role!=nil{
		return role, fmt.Errorf("User already exists")
	}
	newRole:= database.SecurityRole{Name: roleName , Privileges: privileges}
	theNewRole, err := db.AddRole(ctx, newRole); 
	if err!=nil{
		return nil,err
	}
	return &theNewRole, nil
}

func DeleteRole(ctx context.Context, roleName string) error {
	role, err:=db.FindRole(ctx, roleName)
	if err!=nil{
		return err
	}
	if role==nil{
		return fmt.Errorf("Role does not exists")
	}
	users, err:= db.GetAllUsersWithRole(ctx, roleName)
	if err!=nil{
		return err
	}
	if len(users)>0{
		return fmt.Errorf("Users have the role %v cannot delete", roleName)
	}
	return db.RemoveRole(ctx, *role)
}


func UpdateRole(ctx context.Context, roleName string, privileges ...database.SecurityPrivilege) error {
	role, err:=db.FindRole(ctx, roleName)
	if err!=nil{
		return err
	}
	if role==nil{
		return fmt.Errorf("Role does not exists")
	}
	role.Privileges = privileges
	_, err=db.UpdateRole(ctx, *role)
	
	return err
}

func GetRole(ctx context.Context, roleName string) (*database.SecurityRole, error){
	role, err:=db.FindRole(ctx, roleName)
	return role, err
}

func GetAllRoles(ctx context.Context) ([]database.SecurityRole, error){
	return db.GetAllRoles(ctx)
}

func AddRoleToUser(ctx context.Context, user database.UserData ,roleName string)error{
	if err:=user.AddRole(ctx, roleName, db); err!=nil{
		return err
	}
	return nil
}

func IsRoleAuthorized(ctx context.Context, roleName string, privilege database.SecurityPrivilege) (bool,error){
	role, err:=db.FindRole(ctx, roleName)
	if err!=nil{
		return false, err
	}
	if role.HasPrivilege(privilege) || role.HasPrivilege(database.SecurityAdmin) {
		return true, nil
	}
	return false, nil
}

func GetUserRole(ctx context.Context, user database.UserData) (database.SecurityRole, error){
	role, err:=db.FindRole(ctx, user.Role)
	if err!=nil{
		return database.SecurityRole{}, err
	}
	return *role, nil
}

