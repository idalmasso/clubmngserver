package database

import "context"

type ClubDb interface{
	Init()
	FindUser(context.Context, string) (*UserData, error)
	AddUser(context.Context,UserData) (UserData,error)
	UpdateUser(context.Context,UserData) (UserData,error)
	RemoveUser(context.Context,UserData) error
	AddRole(ctx context.Context,role SecurityRole) (SecurityRole,error)
	UpdateRole(ctx context.Context,role SecurityRole) (SecurityRole,error)
	RemoveRole(ctx context.Context,role SecurityRole) error
	FindRole(ctx context.Context,roleName  string) (*SecurityRole, error)
	GetAllRoles(ctx context.Context) ([]SecurityRole, error)
	GetAllUsers(ctx context.Context) ([]UserData, error)
	GetAllUsersWithRole(ctx context.Context, roleName string)([]UserData, error)
}
