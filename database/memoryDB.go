package database

import (
	"context"
	"fmt"
)

type MemoryDB struct{
	users map[string]UserData
	roles map[string]SecurityRole
}


func (db *MemoryDB) Init(){

	if db.roles==nil{
		db.roles=make(map[string]SecurityRole)
		var adminRole SecurityRole
		adminRole.Name="Admin"
		adminRole.Privileges= make([]SecurityPrivilege, 0)
		adminRole.Privileges = append(adminRole.Privileges, SecurityAdmin)
		db.AddRole(context.Background(), adminRole) 
	}
	if db.users==nil{
		db.users=make(map[string]UserData)
		var uData UserData
		uData.Username="administrator"
		uData, err:=db.AddUser(context.Background(), uData)
		if err!=nil{
			panic("Cannot create admin user")
		}
		uData.SetPassword(context.Background(),"Abcd1234", db)
		uData.AddRole(context.Background(), "Admin", db)
	}
}

func (db *MemoryDB) FindUser(ctx context.Context,username  string) (*UserData, error){
	u, ok:=db.users[username]
	if !ok{
		return nil, nil
	}
	return &u, nil
}

func (db *MemoryDB) GetAllUsers(ctx context.Context) ([]UserData, error){
	var uData []UserData
	for _,u:=range(db.users){
		uData=append(uData, u)
	}
	return uData, nil
} 
func (db *MemoryDB) GetAllUsersWithRole(ctx context.Context, roleName string) ([]UserData, error){
	var uData []UserData
	for _,u:=range(db.users){
		if u.Role==roleName{
			uData=append(uData, u)
		}
	}
	return uData, nil
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
func (db *MemoryDB) FindRole(ctx context.Context,roleName  string) (*SecurityRole, error){
	role, ok:=db.roles[roleName]
	if !ok{
		return nil, nil
	}
	return &role, nil
}
func (db *MemoryDB) AddRole(ctx context.Context,role SecurityRole) (SecurityRole,error){
	_, ok:=db.roles[role.Name]
	if ok{
		return SecurityRole{}, fmt.Errorf("Role already exists")
	}
	if role.Privileges ==nil{
		role.Privileges=make([]SecurityPrivilege, 0)
	}
	
	db.roles[role.Name]=role
	
	return db.roles[role.Name],nil
}

func (db *MemoryDB) UpdateRole(ctx context.Context,role SecurityRole) (SecurityRole,error){
	db.roles[role.Name]=role
	return db.roles[role.Name], nil
}

func (db *MemoryDB) RemoveRole(ctx context.Context,role SecurityRole) error{
	delete(db.roles, role.Name)
	return nil
}

func (db *MemoryDB) GetAllRoles(ctx context.Context) ([]SecurityRole, error){
	var roles []SecurityRole
	for _,role:=range(db.roles){
		roles=append(roles, role)
	}
	return roles, nil
}
