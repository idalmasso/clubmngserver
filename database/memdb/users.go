package memdb

import (
	"context"
	"fmt"

	"github.com/idalmasso/clubmngserver/common"
)

//FindUser returns a database.UserData from username
func (db *MemoryDB) FindUser(ctx context.Context,username  string) (*common.UserData, error){
	u, ok:=db.users[username]
	if !ok{
		return nil, nil
	}
	return &u, nil
}
//GetAllUsers returns all users
func (db *MemoryDB) GetAllUsers(ctx context.Context) ([]common.UserData, error){
	var uData []common.UserData
	for _,u:=range(db.users){
		uData=append(uData, u)
	}
	return uData, nil
} 
//GetAllUsersWithRole returns all users with a role
func (db *MemoryDB) GetAllUsersWithRole(ctx context.Context, roleName string) ([]common.UserData, error){
	var uData []common.UserData
	for _,u:=range(db.users){
		if u.Role==roleName{
			uData=append(uData, u)
		}
	}
	return uData, nil
}
//AddUser add a user to the db
func (db *MemoryDB) AddUser(ctx context.Context,user common.UserData) (common.UserData,error){
	_, ok:=db.users[user.Username]
	if ok{
		return common.UserData{}, fmt.Errorf("Username already exists")
	}
	if user.AuthenticationTokens==nil{
		user.AuthenticationTokens=make(map[string]string)
	}
	if user.AuthorizationTokens==nil{
		user.AuthorizationTokens=make(map[string]struct{})
	}
	db.users[user.Username]=user
	
	return db.users[user.Username],nil
}
//UpdateUser update the user in the database. In memory does not throw an error
func (db *MemoryDB) UpdateUser(ctx context.Context,user common.UserData) (common.UserData,error){
	db.users[user.Username]=user
	return db.users[user.Username], nil
}
//RemoveUser remove an actual user from database
func (db *MemoryDB) RemoveUser(ctx context.Context,user common.UserData) error{
	delete(db.users, user.Username)
	return nil
}
