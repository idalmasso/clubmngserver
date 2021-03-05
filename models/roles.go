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
