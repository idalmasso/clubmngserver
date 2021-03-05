package database

import (
	"context"
	"fmt"
)

type MemoryDB struct{
	users map[string]UserData
}


func (db *MemoryDB) Init(){
	if db.users==nil{
		db.users=make(map[string]UserData)
	}
}

func (db *MemoryDB) FindUser(ctx context.Context,username  string) (UserData, error){
	u, ok:=db.users[username]
	if !ok{
		return UserData{}, fmt.Errorf("Not found")
	}
	return u, nil
}

func (db *MemoryDB) AddUser(ctx context.Context,user UserData) (UserData,error){
	_, ok:=db.users[user.Username]
	if ok{
		return UserData{}, fmt.Errorf("Username already exists")
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

func (db *MemoryDB) UpdateUser(ctx context.Context,user UserData) (UserData,error){
	db.users[user.Username]=user
	return db.users[user.Username], nil
}

func (db *MemoryDB) RemoveUser(ctx context.Context,user UserData) error{
	delete(db.users, user.Username)
	return nil
}
