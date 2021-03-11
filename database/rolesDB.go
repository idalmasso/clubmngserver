package database

import (
	"context"

	"github.com/idalmasso/clubmngserver/common"
)

//SecurityRoleDBInterface is the interface a db should implement for the functionality of a database
type SecurityRoleDBInterface interface{
	AddRole(ctx context.Context,role common.SecurityRole) (*common.SecurityRole,error)
	UpdateRole(ctx context.Context,role common.SecurityRole) (*common.SecurityRole,error)
	RemoveRole(ctx context.Context,role common.SecurityRole) error
	FindRole(ctx context.Context,roleName  string) (*common.SecurityRole, error)
	GetAllRoles(ctx context.Context) ([]common.SecurityRole, error)
}
