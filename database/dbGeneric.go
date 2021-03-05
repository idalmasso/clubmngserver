package database

import "context"

type ClubDb interface{
	Init()
	FindUser(context.Context, string) (UserData, error)
	AddUser(context.Context,UserData) (UserData,error)
	UpdateUser(context.Context,UserData) (UserData,error)
	RemoveUser(context.Context,UserData) error

}
