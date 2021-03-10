package models

import (
	"context"
	"fmt"

	"github.com/idalmasso/clubmngserver/common"
)

func AddRole(ctx context.Context, roleName string, privileges ...common.SecurityPrivilege) ( *common.SecurityRole,error) {
	role, err:=db.FindRole(ctx, roleName)
	if err!=nil{
		return nil,err
	}
	if role!=nil{
		return role, fmt.Errorf("User already exists")
	}
	newRole:= common.SecurityRole{Name: roleName , Privileges: privileges}
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


func UpdateRole(ctx context.Context, roleName string, privileges ...common.SecurityPrivilege) error {
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

func GetRole(ctx context.Context, roleName string) (*common.SecurityRole, error){
	role, err:=db.FindRole(ctx, roleName)
	return role, err
}

func GetAllRoles(ctx context.Context) ([]common.SecurityRole, error){
	return db.GetAllRoles(ctx)
}

func AddRoleToUser(ctx context.Context, user common.UserData ,roleName string)error{
	role, err:=db.FindRole(ctx, roleName)
	if err!=nil{
		return err
	}
	user.AddRole(ctx, *role)
	if _, err:=db.UpdateUser(ctx, user); err!=nil{
		return err
	}
	return nil
}

func IsRoleAuthorized(ctx context.Context, roleName string, privilege common.SecurityPrivilege) (bool,error){
	role, err:=db.FindRole(ctx, roleName)
	if err!=nil{
		return false, err
	}
	if role.HasPrivilege(privilege) || role.HasPrivilege(common.SecurityAdmin) {
		return true, nil
	}
	return false, nil
}

func GetUserRole(ctx context.Context, user common.UserData) (common.SecurityRole, error){
	role, err:=db.FindRole(ctx, user.Role)
	if err!=nil{
		return common.SecurityRole{}, err
	}
	return *role, nil
}

