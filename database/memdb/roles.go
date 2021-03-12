package memdb

import (
	"context"

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
		return nil, common.AlreadyExistsError{ID:role.Name}
	}
	if role.Privileges ==nil{
		role.Privileges=make([]common.SecurityPrivilege, 0)
	}
	
	db.roles[role.Name]=&role
	
	return db.roles[role.Name],nil
}
//UpdateRole updates an actual role
func (db *MemoryDB) UpdateRole(ctx context.Context,role common.SecurityRole) (*common.SecurityRole,error){
	if _, ok:=db.roles[role.Name]; !ok{
		return nil, common.NotFoundError{ID:role.Name}
	}
	db.roles[role.Name]=&role
	return db.roles[role.Name], nil
}
//RemoveRole remove a role from the database
func (db *MemoryDB) RemoveRole(ctx context.Context,roleName string) error{
	delete(db.roles, roleName)
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
