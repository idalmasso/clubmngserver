package models

import (
	"context"

	"github.com/idalmasso/clubmngserver/database"
)



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
	return role, nil
}
