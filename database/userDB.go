package database

import (
	"context"

	"github.com/idalmasso/clubmngserver/common"
)

//UserAuthDBInterface interface for a database to implement to users mng and auth
type UserAuthDBInterface interface{
	FindUser(context.Context, string) (*common.UserData, error)
	AddUser(context.Context,common.UserData) (*common.UserData,error)
	UpdateUser(context.Context,common.UserData) (*common.UserData,error)
	RemoveUser(context.Context,common.UserData) error
	GetAllUsers( context.Context) ([]common.UserData, error)
	GetAllUsersWithRole( context.Context,  string)([]common.UserData, error)	
}

