package database

import (
	"context"

	"github.com/idalmasso/clubmngserver/common"
)

//SecurityRoleDBInterface is the interface a db should implement for the functionality of a database
type SecurityRoleDBInterface interface{
	AddRole( context.Context, common.SecurityRole) (*common.SecurityRole,error)
	UpdateRole( context.Context, common.SecurityRole) (*common.SecurityRole,error)
	RemoveRole( context.Context, string) error
	FindRole( context.Context,  string) (*common.SecurityRole, error)
	GetAllRoles( context.Context) ([]common.SecurityRole, error)
}
