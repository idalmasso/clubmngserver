package database

import (
	"context"

	"github.com/idalmasso/clubmngserver/common/userprops"
)

//UserPropsDB interface for a database to implement to users properties
type UserPropsDB interface{
	FindUserPropertyDefinition(context.Context, string) (userprops.UserPropertyDefinition, error)
	AddUserPropertyDefinition(context.Context,userprops.UserPropertyDefinition) (userprops.UserPropertyDefinition,error)
	UpdateUserPropertyDefinition(context.Context,userprops.UserPropertyDefinition) (userprops.UserPropertyDefinition,error)
	RemoveUserPropertyDefinition( context.Context, string) error
	GetAllUserPropertyDefinitions(context.Context) ([]userprops.UserPropertyDefinition, error)
}

