package memdb

import (
	"context"
	"fmt"

	"github.com/idalmasso/clubmngserver/common"
)

//FindRole find and returns a Role from name
func (db *MemoryDB) FindRole(ctx context.Context,roleName  string) (*common.SecurityRole, error){
	role, ok:=db.roles[roleName]
	if !ok{
		return nil, nil
	}
	return role, nil
}

//AddRole add a new empty role to the database
func (db *MemoryDB) AddRole(ctx context.Context,role common.SecurityRole) (*common.SecurityRole,error){
	_, ok:=db.roles[role.Name]
	if ok{
		return nil, fmt.Errorf("Role already exists")
	}
	if role.Privileges ==nil{
		role.Privileges=make([]common.SecurityPrivilege, 0)
	}
	
	db.roles[role.Name]=&role
	
	return db.roles[role.Name],nil
}
//UpdateRole updates an actual role
func (db *MemoryDB) UpdateRole(ctx context.Context,role common.SecurityRole) (*common.SecurityRole,error){
	db.roles[role.Name]=&role
	return db.roles[role.Name], nil
}
//RemoveRole remove a role from the database
func (db *MemoryDB) RemoveRole(ctx context.Context,role common.SecurityRole) error{
	delete(db.roles, role.Name)
	return nil
}
//GetAllRoles returns all roles in the db
func (db *MemoryDB) GetAllRoles(ctx context.Context) ([]common.SecurityRole, error){
	var roles []common.SecurityRole
	for _,role:=range(db.roles){
		roles=append(roles, *role)
	}
	return roles, nil
}
